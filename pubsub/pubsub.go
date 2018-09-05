// Copyright (c) 2017,2018 Zededa, Inc.
// All rights reserved.

// Provide for a pubsub mechanism for config and status which is
// backed by an IPC mechanism such as connected sockets.

package pubsub

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/zededa/go-provision/watch"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Protocol over AF_UNIX or other IPC mechanism
// XXX add receive side code later
// "request" from client after connect to sanity check subject.
// Server sends the other messages; "update" for initial values.
// "complete" once all initial keys/values in collection have been sent.
// "restarted" if/when pub.km.restarted is set.
// Ongoing we send "update" and "delete" messages.
// They keys and values are base64-encoded since they might contain spaces.
// We include typeName after command word for sanity checks.
// Hence the message format is
//	"request" typeName
//	"hello"  string
//	"update" typeName key json-val
//	"delete" typeName key
//	"complete" typeName
//	"restarted" typeName

// Maintain a collection which is used to handle the restart of a subscriber
// map of agentname, key to get a json string
// We use sync.Map to allow concurrent access, and handle notifications
// to work with the relaxed output from the sync.Map's Range function
type keyMap struct {
	restarted bool
	key       sync.Map
}

// We always publish to our collection.
// We always write to a file in order to have a checkpoint on restart
// The special agent name "" implies always reading from the /var/run/zededa/
// directory.
const publishToSock = true     // XXX
const subscribeFromDir = false // XXX
const subscribeFromSock = true // XXX

// For a subscription, if the agentName is empty we interpret that as
// being directory in /var/tmp/zededa
const fixedName = "zededa"
const fixedDir = "/var/tmp/" + fixedName

const debug = true // XXX setable?

type notify struct{}

// The set of channels to which we need to send notifications
type updaters struct {
	lock    sync.Mutex
	servers []chan<- notify
}

var updaterList updaters

func updatersAdd(updater chan notify) {
	updaterList.lock.Lock()
	updaterList.servers = append(updaterList.servers, updater)
	updaterList.lock.Unlock()
}

func updatersRemove(updater chan notify) {
	updaterList.lock.Lock()
	servers := make([]chan<- notify, len(updaterList.servers))
	found := false
	for _, old := range updaterList.servers {
		if old == updater {
			found = true
		} else {
			servers = append(servers, old)
		}
	}
	if !found {
		log.Fatal("updatersRemove: not found\n")
	}
	updaterList.servers = servers
	updaterList.lock.Unlock()
}

// Send a notification to all the channels which does not yet
// have one queued
func updatersNotify() {
	updaterList.lock.Lock()
	for i, server := range updaterList.servers {
		select {
		case server <- notify{}:
			log.Printf("updaterNotify sent to %d\n", i)
		default:
			log.Printf("updaterNotify NOT sent to %d\n", i)
		}
	}
	updaterList.lock.Unlock()
}

// Usage:
//  p1, err := pubsub.Publish("foo", fooStruct{})
//  ...
//  // Optional
//  p1.SignalRestarted()
//  ...
//  p1.Publish(key, item)
//  p1.Unpublish(key) to delete
//
//  foo := p1.Get(key)
//  fooAll := p1.GetAll()

type Publication struct {
	// Private fields
	topicType  interface{}
	agentName  string
	agentScope string
	topic      string
	km         keyMap
	sockName   string
	listener   net.Listener
}

func Publish(agentName string, topicType interface{}) (*Publication, error) {
	return publishImpl(agentName, "", topicType)
}

func PublishScope(agentName string, agentScope string, topicType interface{}) (*Publication, error) {
	return publishImpl(agentName, agentScope, topicType)
}

// Init function to create directory and socket listener based on above settings
// We read any checkpointed state from dirName and insert in pub.km as initial
// values.
func publishImpl(agentName string, agentScope string,
	topicType interface{}) (*Publication, error) {

	topic := TypeToName(topicType)
	pub := new(Publication)
	pub.topicType = topicType
	pub.agentName = agentName
	pub.agentScope = agentScope
	pub.topic = topic
	name := pub.nameString()

	log.Printf("Publish(%s)\n", name)

	// We always write to the directory as a checkpoint
	dirName := PubDirName(name)
	if _, err := os.Stat(dirName); err != nil {
		log.Printf("Publish Create %s\n", dirName)
		if err := os.MkdirAll(dirName, 0700); err != nil {
			errStr := fmt.Sprintf("Publish(%s): %s",
				name, err)
			return nil, errors.New(errStr)
		}
	} else {
		// Read existig status from dir
		pub.populate()
		if debug {
			pub.dump("after populate")
		}
	}

	if publishToSock {
		sockName := SockName(name)
		if _, err := os.Stat(sockName); err == nil {
			if err := os.Remove(sockName); err != nil {
				errStr := fmt.Sprintf("Publish(%s): %s",
					name, err)
				return nil, errors.New(errStr)
			}
		}
		s, err := net.Listen("unixpacket", sockName)
		if err != nil {
			errStr := fmt.Sprintf("Publish(%s): failed %s",
				name, err)
			return nil, errors.New(errStr)
		}
		pub.sockName = sockName
		pub.listener = s
		go pub.publisher()
	}
	return pub, nil
}

// Only reads json files. Sets restarted if that file was found.
func (pub *Publication) populate() {
	name := pub.nameString()
	dirName := PubDirName(name)
	foundRestarted := false

	log.Printf("populate(%s)\n", name)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			if file.Name() == "restarted" {
				foundRestarted = true
			}
			continue
		}
		// Remove .json from name */
		key := strings.Split(file.Name(), ".json")[0]

		statusFile := dirName + "/" + file.Name()
		if _, err := os.Stat(statusFile); err != nil {
			// File just vanished!
			log.Printf("populate: File disappeared <%s>\n",
				statusFile)
			continue
		}

		log.Printf("populate found key %s file %s\n", key, statusFile)

		sb, err := ioutil.ReadFile(statusFile)
		if err != nil {
			log.Printf("populate: %s for %s\n", err, statusFile)
			continue
		}
		var item interface{}
		if err := json.Unmarshal(sb, &item); err != nil {
			log.Printf("populate: %s file: %s\n",
				err, statusFile)
			continue
		}
		pub.km.key.Store(key, item)
	}
	pub.km.restarted = foundRestarted
	log.Printf("populate(%s) done\n", name)
}

// go routine which runs the AF_UNIX server.
func (pub *Publication) publisher() {
	name := pub.nameString()
	for {
		c, err := pub.listener.Accept()
		if err != nil {
			log.Printf("publisher(%s) failed %s\n", name, err)
			continue
		}
		go pub.serveConnection(c)
	}
}

// Used locally by each serverConnection goroutine to track updates
// to send.
type localCollection map[string]interface{}

func (pub *Publication) serveConnection(s net.Conn) {
	name := pub.nameString()
	log.Printf("serveConnection(%s)\n", name)
	defer s.Close()

	// Track the set of keys/values we are sending to the peer
	sendToPeer := make(localCollection)
	sentRestarted := false
	// Read request
	buf := make([]byte, 65536)
	res, err := s.Read(buf)
	// XXX check if res == 65536 i.e., truncated
	request := strings.Split(string(buf[0:res]), " ")
	log.Printf("serveConnection read %d: %v\n", len(request), request)
	if len(request) != 2 || request[0] != "request" || request[1] != pub.topic {
		log.Printf("Invalid request message: %v\n", request)
		return
	}

	_, err = s.Write([]byte(fmt.Sprintf("hello %s", name)))
	if err != nil {
		log.Printf("serveConnection(%s) failed %s\n",
			name, err)
		return
	}
	// Insert our notification channel before we get the initial
	// snapshot to avoid missing any updates/deletes.
	updater := make(chan notify)
	updatersAdd(updater)
	defer updatersRemove(updater)

	// Get a local snapshot of the collection and the set of keys
	// we need to send these. Updates the slave collection.
	keys := pub.determineDiffs(sendToPeer)

	// Send the keys we just determined; all since this is the initial
	err = pub.serialize(s, keys, sendToPeer)
	if err != nil {
		log.Printf("serveConnection(%s) failed %s\n",
			name, err)
		return
	}
	err = pub.sendComplete(s)
	if err != nil {
		log.Printf("serveConnection(%s) failed %s\n",
			name, err)
		return
	}
	if pub.km.restarted && !sentRestarted {
		err = pub.sendRestarted(s)
		if err != nil {
			log.Printf("serveConnection(%s) failed %s\n",
				name, err)
			return
		}
		sentRestarted = true
	}

	// Handle any changes
	for {
		<-updater
		log.Printf("Received notification\n")
		// Update and determine which keys changed
		keys := pub.determineDiffs(sendToPeer)

		// Send the updates and deletes for those keys
		err = pub.serialize(s, keys, sendToPeer)
		if err != nil {
			log.Printf("serveConnection(%s) failed %s\n",
				name, err)
			return
		}

		if pub.km.restarted && !sentRestarted {
			err = pub.sendRestarted(s)
			if err != nil {
				log.Printf("serveConnection(%s) failed %s\n",
					name, err)
				return
			}
			sentRestarted = true
		}
	}
}

// Returns the deleted keys before the added/modified ones
func (pub *Publication) determineDiffs(slaveCollection localCollection) []string {
	var keys []string
	items := pub.GetAll()
	// Look for deleted
	for slaveKey, _ := range slaveCollection {
		master, _ := pub.Get(slaveKey)
		if master == nil {
			log.Printf("determineDiffs: key %s deleted\n",
				slaveKey)
			delete(slaveCollection, slaveKey)
			keys = append(keys, slaveKey)
		}
	}
	// Look for new/changed
	for masterKey, master := range items {
		slave := lookupSlave(slaveCollection, masterKey)
		if slave == nil {
			log.Printf("determineDiffs: key %s added\n",
				masterKey)
			slaveCollection[masterKey] = master
			keys = append(keys, masterKey)
		} else if !cmp.Equal(master, slave) {
			log.Printf("determineDiffs: key %s changed %v\n",
				masterKey, cmp.Diff(master, slave))
			slaveCollection[masterKey] = master
			keys = append(keys, masterKey)
		} else {
			log.Printf("determineDiffs: key %s unchanged\n",
				masterKey)
		}
	}
	return keys
}

func lookupSlave(slaveCollection localCollection, key string) *interface{} {
	for slaveKey, _ := range slaveCollection {
		if slaveKey == key {
			res := slaveCollection[slaveKey]
			return &res
		}
	}
	return nil
}

func TypeToName(something interface{}) string {
	t := reflect.TypeOf(something)
	out := strings.Split(t.String(), ".")
	return out[len(out)-1]
}

func SockName(name string) string {
	return fmt.Sprintf("/var/run/%s.sock", name)
}

func PubDirName(name string) string {
	return fmt.Sprintf("/var/run/%s", name)
}

func FixedDirName(name string) string {
	return fmt.Sprintf("%s/%s", fixedDir, name)
}

func (pub *Publication) nameString() string {
	if pub.agentScope == "" {
		return fmt.Sprintf("%s/%s", pub.agentName, pub.topic)
	} else {
		return fmt.Sprintf("%s/%s/%s", pub.agentName, pub.agentScope,
			pub.topic)
	}
}

func (pub *Publication) Publish(key string, item interface{}) error {
	topic := TypeToName(item)
	name := pub.nameString()
	if topic != pub.topic {
		errStr := fmt.Sprintf("Publish(%s): item is wrong topic %s",
			name, topic)
		return errors.New(errStr)
	}
	// Perform a deepCopy so the Equal check will work
	newItem := deepCopy(item)
	if m, ok := pub.km.key.Load(key); ok {
		if cmp.Equal(m, newItem) {
			if debug {
				log.Printf("Publish(%s/%s) unchanged\n",
					name, key)
			}
			return nil
		}
		if debug {
			log.Printf("Publish(%s/%s) replacing due to diff %s\n",
				name, key, cmp.Diff(m, newItem))
		}
	} else if debug {
		log.Printf("Publish(%s/%s) adding %+v\n",
			name, key, newItem)
	}
	pub.km.key.Store(key, newItem)

	if debug {
		pub.dump("after Publish")
	}
	dirName := PubDirName(name)
	fileName := dirName + "/" + key + ".json"
	if debug {
		log.Printf("Publish writing %s\n", fileName)
	}
	// XXX already did a marshal in deepCopy; save that result?
	b, err := json.Marshal(item)
	if err != nil {
		log.Fatal(err, "json Marshal in Publish")
	}
	err = WriteRename(fileName, b)
	if err != nil {
		return err
	}
	updatersNotify()
	return nil
}

func WriteRename(fileName string, b []byte) error {
	dirName := filepath.Dir(fileName)
	// Do atomic rename to avoid partially written files
	tmpfile, err := ioutil.TempFile(dirName, "pubsub")
	if err != nil {
		errStr := fmt.Sprintf("WriteRename(%s): %s",
			fileName, err)
		return errors.New(errStr)
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.Write(b)
	if err != nil {
		errStr := fmt.Sprintf("WriteRename(%s): %s",
			fileName, err)
		return errors.New(errStr)
	}
	if err := tmpfile.Close(); err != nil {
		errStr := fmt.Sprintf("WriteRename(%s): %s",
			fileName, err)
		return errors.New(errStr)
	}
	if err := os.Rename(tmpfile.Name(), fileName); err != nil {
		errStr := fmt.Sprintf("WriteRename(%s): %s",
			fileName, err)
		return errors.New(errStr)
	}
	return nil
}

func deepCopy(in interface{}) interface{} {
	b, err := json.Marshal(in)
	if err != nil {
		log.Fatal(err, "json Marshal in deepCopy")
	}
	var output interface{}
	if err := json.Unmarshal(b, &output); err != nil {
		log.Fatal(err, "json Unmarshal in deepCopy")
	}
	return output
}

func (pub *Publication) Unpublish(key string) error {
	name := pub.nameString()
	if m, ok := pub.km.key.Load(key); ok {
		if debug {
			log.Printf("Unpublish(%s/%s) removing %+v\n",
				name, key, m)
		}
	} else {
		errStr := fmt.Sprintf("Unpublish(%s/%s): key does not exist",
			name, key)
		log.Printf("%s\n", errStr)
		return errors.New(errStr)
	}
	pub.km.key.Delete(key)
	if debug {
		pub.dump("after Unpublish")
	}
	updatersNotify()

	dirName := PubDirName(name)
	fileName := dirName + "/" + key + ".json"
	if debug {
		log.Printf("Unpublish deleting file %s\n", fileName)
	}
	if err := os.Remove(fileName); err != nil {
		errStr := fmt.Sprintf("Unpublish(%s/%s): failed %s",
			name, key, err)
		return errors.New(errStr)
	}
	return nil
}

func (pub *Publication) SignalRestarted() error {
	if debug {
		log.Printf("pub.SignalRestarted(%s)\n", pub.nameString())
	}
	return pub.restartImpl(true)
}

func (pub *Publication) ClearRestarted() error {
	if debug {
		log.Printf("pub.ClearRestarted(%s)\n", pub.nameString())
	}
	return pub.restartImpl(false)
}

// Record the restarted state and send over socket/file.
func (pub *Publication) restartImpl(restarted bool) error {
	name := pub.nameString()
	log.Printf("pub.restartImpl(%s, %v)\n", name, restarted)
	if restarted == pub.km.restarted {
		log.Printf("pub.restartImpl(%s, %v) value unchanged\n",
			name, restarted)
		return nil
	}
	pub.km.restarted = restarted
	if restarted {
		updatersNotify()
	}

	dirName := PubDirName(name)
	restartFile := dirName + "/" + "restarted"
	if restarted {
		f, err := os.OpenFile(restartFile, os.O_RDONLY|os.O_CREATE, 0600)
		if err != nil {
			errStr := fmt.Sprintf("pub.restartImpl(%s): openfile failed %s",
				name, err)
			return errors.New(errStr)
		}
		f.Close()
	} else {
		if err := os.Remove(restartFile); err != nil {
			errStr := fmt.Sprintf("pub.restartImpl(%s): remove failed %s",
				name, err)
			return errors.New(errStr)
		}
	}
	return nil
}

func (pub *Publication) serialize(sock net.Conn, keys []string,
	sendToPeer localCollection) error {

	name := pub.nameString()
	log.Printf("serialize(%s, %v)\n", name, keys)
	// Any initial deletes?
	for _, key := range keys {
		val, ok := sendToPeer[key]
		if ok {
			err := pub.sendUpdate(sock, key, val)
			if err != nil {
				log.Printf("serialize(%s) write failed %s\n",
					name, err)
				return err
			}
		} else {
			err := pub.sendDelete(sock, key)
			if err != nil {
				log.Printf("serialize(%s) write failed %s\n",
					name, err)
				return err
			}
		}
	}
	return nil
}

func (pub *Publication) sendUpdate(sock net.Conn, key string,
	val interface{}) error {

	if debug {
		log.Printf("sendUpdate: key %s\n", key)
	}
	b, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err, "json Marshal in serialize")
	}
	// base64-encode to avoid having spaces in the key and val
	sendKey := base64.StdEncoding.EncodeToString([]byte(key))
	sendVal := base64.StdEncoding.EncodeToString(b)
	_, err = sock.Write([]byte(fmt.Sprintf("update %s %s %s",
		pub.topic, sendKey, sendVal)))
	return err
}

func (pub *Publication) sendDelete(sock net.Conn, key string) error {
	if debug {
		log.Printf("sendDelete: key %s\n", key)
	}
	// base64-encode to avoid having spaces in the key
	sendKey := base64.StdEncoding.EncodeToString([]byte(key))
	_, err := sock.Write([]byte(fmt.Sprintf("delete %s %s",
		pub.topic, sendKey)))
	return err
}

func (pub *Publication) sendRestarted(sock net.Conn) error {
	if debug {
		log.Printf("sendRestarted\n")
	}
	_, err := sock.Write([]byte(fmt.Sprintf("restarted %s", pub.topic)))
	return err
}

func (pub *Publication) sendComplete(sock net.Conn) error {
	if debug {
		log.Printf("sendComplete\n")
	}
	_, err := sock.Write([]byte(fmt.Sprintf("complete %s", pub.topic)))
	return err
}

func (pub *Publication) dump(infoStr string) {
	name := pub.nameString()
	log.Printf("dump(%s) %s\n", name, infoStr)
	dumper := func(key, val interface{}) bool {
		b, err := json.Marshal(val)
		if err != nil {
			log.Fatal(err, "json Marshal in dump")
		}
		log.Printf("\tkey %s val %s\n", key.(string), b)
		return true
	}
	pub.km.key.Range(dumper)

	log.Printf("\trestarted %t\n", pub.km.restarted)
}

func (pub *Publication) Get(key string) (interface{}, error) {
	m, ok := pub.km.key.Load(key)
	if ok {
		return m, nil
	} else {
		name := pub.nameString()
		errStr := fmt.Sprintf("Get(%s) unknown key %s", name, key)
		return nil, errors.New(errStr)
	}
}

// Enumerate all the key, value for the collection
func (pub *Publication) GetAll() map[string]interface{} {
	result := make(map[string]interface{})
	assigner := func(key, val interface{}) bool {
		result[key.(string)] = val
		return true
	}
	pub.km.key.Range(assigner)
	return result
}

// Usage:
//  s1 := pubsub.Subscribe("foo", fooStruct{}, true, &myctx)
// Or
//  s1 := pubsub.Subscribe("foo", fooStruct{}, false, &myctx)
//  s1.ModifyHandler = func(...), // Optional
//  s1.DeleteHandler = func(...), // Optional
//  s1.RestartHandler = func(...), // Optional
//  [ Initialize myctx ]
//  s1.Activate()
//  ...
//  select {
//     change := <- s1.C:
//         s1.ProcessChange(change, ctx)
//  }
//  The ProcessChange function calls the various handlers (if set) and updates
//  the subscribed collection. The subscribed collection can be accessed using:
//  foo := s1.Get(key)
//  fooAll := s1.GetAll()

type SubModifyHandler func(ctx interface{}, key string, status interface{})
type SubDeleteHandler func(ctx interface{}, key string, status interface{})
type SubRestartHandler func(ctx interface{}, restarted bool)

type Subscription struct {
	C              <-chan string
	ModifyHandler  SubModifyHandler
	DeleteHandler  SubDeleteHandler
	RestartHandler SubRestartHandler

	// Private fields
	sendChan   chan<- string
	topicType  interface{}
	agentName  string
	agentScope string
	topic      string
	km         keyMap
	userCtx    interface{}
	sock       net.Conn // For socket subscriptions
	// Handle special case of file only info
	subscribeFromDir bool
	dirName          string
}

func (sub *Subscription) nameString() string {
	agentName := sub.agentName
	if agentName == "" {
		agentName = fixedName
	}
	if sub.agentScope == "" {
		return fmt.Sprintf("%s/%s", sub.agentName, sub.topic)
	} else {
		return fmt.Sprintf("%s/%s/%s", sub.agentName, sub.agentScope,
			sub.topic)
	}
}

// Init function for Subscribe; returns a context.
// Assumption is that agent with call Get(key) later or specify
// handleModify and/or handleDelete functions
// watch ensures that any restart/restarted notification is after any other
// notifications from ReadDir
func Subscribe(agentName string, topicType interface{}, activate bool,
	ctx interface{}) (*Subscription, error) {

	return subscribeImpl(agentName, "", topicType, activate, ctx)
}

func SubscribeScope(agentName string, agentScope string, topicType interface{},
	activate bool, ctx interface{}) (*Subscription, error) {

	return subscribeImpl(agentName, agentScope, topicType, activate, ctx)
}

func subscribeImpl(agentName string, agentScope string, topicType interface{},
	activate bool, ctx interface{}) (*Subscription, error) {

	topic := TypeToName(topicType)
	changes := make(chan string)
	sub := new(Subscription)
	sub.C = changes
	sub.sendChan = changes
	sub.topicType = topicType
	sub.agentName = agentName
	sub.agentScope = agentScope
	sub.topic = topic
	sub.userCtx = ctx
	name := sub.nameString()

	if agentName == "" {
		sub.subscribeFromDir = true
		sub.dirName = FixedDirName(name)
	} else {
		sub.subscribeFromDir = subscribeFromDir
		sub.dirName = PubDirName(name)
	}
	log.Printf("Subscribe(%s)\n", name)

	if activate {
		if err := sub.Activate(); err != nil {
			return nil, err
		}
	}
	return sub, nil
}

// If the agentName is empty we interpret that as being dir /var/tmp/zededa
func (sub *Subscription) Activate() error {

	name := sub.nameString()
	if sub.subscribeFromDir {
		// Waiting for directory to appear
		for {
			if _, err := os.Stat(sub.dirName); err != nil {
				errStr := fmt.Sprintf("Subscribe(%s): failed %s; waiting",
					name, err)
				log.Println(errStr)
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}
		go watch.WatchStatus(sub.dirName, sub.sendChan)
		return nil
	} else if subscribeFromSock {
		go sub.watchSock()
		return nil
	} else {
		errStr := fmt.Sprintf("Subscribe(%s): failed %s",
			name, "nowhere to subscribe")
		return errors.New(errStr)
	}
}

func (sub *Subscription) watchSock() {

	for {
		msg, key, val := sub.connectAndRead()
		switch msg {
		case "hello", "complete":
			if debug {
				log.Printf("watchSock: %s\n", msg)
			}
			// XXX anything for complete? Do we have an initial loop?
		case "restarted":
			if debug {
				log.Printf("watchSock: %s\n", msg)
			}
			sub.sendChan <- "R done"

		case "delete":
			if debug {
				log.Printf("delete base64 key %s\n", key)
			}
			sub.sendChan <- "D " + key

		case "update":
			if debug {
				log.Printf("update base64 key %s\n", key)
			}
			// XXX is size of val any issue? pointer?
			sub.sendChan <- "M " + key + " " + val
		}
	}
}

// Returns msg, key, val
// key and val are base64-encoded
func (sub *Subscription) connectAndRead() (string, string, string) {

	name := sub.nameString()
	sockName := SockName(name)
	buf := make([]byte, 65535)

	// Waiting for publisher to appear; retry on error
	for {
		if sub.sock == nil {
			s, err := net.Dial("unixpacket", sockName)
			if err != nil {
				errStr := fmt.Sprintf("connectAndRead(%s): Dial failed %s",
					name, err)
				log.Println(errStr)
				time.Sleep(10 * time.Second)
				continue
			}
			sub.sock = s
			req := fmt.Sprintf("request %s", sub.topic)
			_, err = s.Write([]byte(req))
			if err != nil {
				errStr := fmt.Sprintf("connectAndRead(%s): sock write failed %s",
					name, err)
				log.Println(errStr)
				sub.sock.Close()
				sub.sock = nil
				continue
			}
		}

		res, err := sub.sock.Read(buf)
		if err != nil {
			errStr := fmt.Sprintf("connectAndRead(%s): sock read failed %s",
				name, err)
			log.Println(errStr)
			sub.sock.Close()
			sub.sock = nil
			continue
		}

		// XXX check if res == 65536 i.e. truncated?
		reply := strings.Split(string(buf[0:res]), " ")
		count := len(reply)
		if count < 2 {
			errStr := fmt.Sprintf("connectAndRead(%s): too short read",
				name)
			log.Println(errStr)
			continue
		}
		msg := reply[0]
		t := reply[1]

		// XXX check type against sub.topic

		// XXX are there error cases where we should Close and
		// continue aka reconnect?
		switch msg {
		case "hello", "restarted", "complete":
			if debug {
				log.Printf("connectAndRead(%s) Got message %s type %s\n",
					msg, t)
			}
			return msg, "", ""

		case "delete":
			if count < 3 {
				errStr := fmt.Sprintf("connectAndRead(%s): too short delete",
					name)
				log.Println(errStr)
				continue
			}
			recvKey := reply[2]

			if debug {
				key, err := base64.StdEncoding.DecodeString(recvKey)
				if err != nil {
					errStr := fmt.Sprintf("connectAndRead(%s): base64 failed %s",
						name, err)
					log.Println(errStr)
					continue
				}
				log.Printf("delete type %s key %s\n",
					t, string(key))
			}
			return msg, recvKey, ""

		case "update":
			if count < 4 {
				errStr := fmt.Sprintf("connectAndRead(%s): too short update",
					name)
				log.Println(errStr)
				continue
			}
			if count > 4 {
				errStr := fmt.Sprintf("connectAndRead(%s): too long update",
					name)
				log.Println(errStr)
				continue
			}
			recvKey := reply[2]
			recvVal := reply[3]
			if debug {
				key, err := base64.StdEncoding.DecodeString(recvKey)
				if err != nil {
					errStr := fmt.Sprintf("connectAndRead(%s): base64 failed %s",
						name, err)
					log.Println(errStr)
					continue
				}
				val, err := base64.StdEncoding.DecodeString(recvVal)
				if err != nil {
					errStr := fmt.Sprintf("connectAndRead(%s): base64 val failed %s",
						name, err)
					log.Println(errStr)
					continue
				}
				log.Printf("update type %s key %s val %s\n",
					t, string(key), string(val))
			}
			return msg, recvKey, recvVal

		default:
			errStr := fmt.Sprintf("connectAndRead(%s): unknown message %s",
				name, msg)
			log.Println(errStr)
			continue
		}
	}
}

// XXX note that change filename includes .json for files. Removed by
// HandleStatusEvent
func (sub *Subscription) ProcessChange(change string) {
	name := sub.nameString()
	if debug {
		log.Printf("ProcessChange(%s) %s\n", name, change)
	}
	if sub.subscribeFromDir {
		var restartFn watch.StatusRestartHandler = handleRestart
		watch.HandleStatusEvent(change, sub,
			sub.dirName, &sub.topicType,
			handleModify, handleDelete, &restartFn)
	} else if subscribeFromSock {
		reply := strings.Split(change, " ")
		operation := reply[0]

		switch operation {
		case "R":
			handleRestart(sub, true)
		case "D":
			recvKey := reply[1]
			key, err := base64.StdEncoding.DecodeString(recvKey)
			if err != nil {
				errStr := fmt.Sprintf("ProcessChange(%s): base64 failed %s",
					name, err)
				log.Println(errStr)
				return
			}
			handleDelete(sub, string(key))

		case "M":
			recvKey := reply[1]
			recvVal := reply[2]
			key, err := base64.StdEncoding.DecodeString(recvKey)
			if err != nil {
				errStr := fmt.Sprintf("ProcessChange(%s): base64 failed %s",
					name, err)
				log.Println(errStr)
				return
			}
			val, err := base64.StdEncoding.DecodeString(recvVal)
			if err != nil {
				errStr := fmt.Sprintf("ProcessChange(%s): base64 val failed %s",
					name, err)
				log.Println(errStr)
				return
			}
			var output interface{}
			if err := json.Unmarshal(val, &output); err != nil {
				errStr := fmt.Sprintf("ProcessChange(%s): json failed %s",
					name, err)
				log.Println(errStr)
				return
			}
			handleModify(sub, string(key), output)
		}
	}
}

func handleModify(ctxArg interface{}, key string, item interface{}) {
	sub := ctxArg.(*Subscription)
	name := sub.nameString()
	if debug {
		log.Printf("pubsub.handleModify(%s) key %s\n", name, key)
	}
	// NOTE: without a deepCopy we would just save a pointer since
	// item is a pointer. That would cause failures.
	newItem := deepCopy(item)
	m, ok := sub.km.key.Load(key)
	if ok {
		if cmp.Equal(m, newItem) {
			if debug {
				log.Printf("pubsub.handleModify(%s/%s) unchanged\n",
					name, key)
			}
			return
		}
		if debug {
			log.Printf("pubsub.handleModify(%s/%s) replacing due to diff %s\n",
				name, key, cmp.Diff(m, newItem))
		}
	} else if debug {
		log.Printf("pubsub.handleModify(%s) add %+v for key %s\n",
			name, newItem, key)
	}
	sub.km.key.Store(key, newItem)
	if debug {
		sub.dump("after handleModify")
	}
	if sub.ModifyHandler != nil {
		(sub.ModifyHandler)(sub.userCtx, key, newItem)
	}
	if debug {
		log.Printf("pubsub.handleModify(%s) done for key %s\n",
			name, key)
	}
}

func handleDelete(ctxArg interface{}, key string) {
	sub := ctxArg.(*Subscription)
	name := sub.nameString()
	if debug {
		log.Printf("pubsub.handleDelete(%s) key %s\n", name, key)
	}
	m, ok := sub.km.key.Load(key)
	if !ok {
		log.Printf("pubsub.handleDelete(%s) %s key not found\n",
			name, key)
		return
	}
	if debug {
		log.Printf("pubsub.handleDelete(%s) key %s value %+v\n",
			name, key, m)
	}
	sub.km.key.Delete(key)
	if debug {
		sub.dump("after handleDelete")
	}
	if sub.DeleteHandler != nil {
		(sub.DeleteHandler)(sub.userCtx, key, m)
	}
	if debug {
		log.Printf("pubsub.handleModify(%s) done for key %s\n",
			name, key)
	}
}

func handleRestart(ctxArg interface{}, restarted bool) {
	sub := ctxArg.(*Subscription)
	name := sub.nameString()
	if debug {
		log.Printf("pubsub.handleRestart(%s) restarted %v\n",
			name, restarted)
	}
	if restarted == sub.km.restarted {
		if debug {
			log.Printf("pubsub.handleRestart(%s) value unchanged\n",
				name)
		}
		return
	}
	sub.km.restarted = restarted
	if sub.RestartHandler != nil {
		(sub.RestartHandler)(sub.userCtx, restarted)
	}
	if debug {
		log.Printf("pubsub.handleRestart(%s) done for restarted %v\n",
			name, restarted)
	}
}

func (sub *Subscription) dump(infoStr string) {
	name := sub.nameString()
	log.Printf("dump(%s) %s\n", name, infoStr)
	dumper := func(key, val interface{}) bool {
		b, err := json.Marshal(val)
		if err != nil {
			log.Fatal(err, "json Marshal in dump")
		}
		log.Printf("\tkey %s val %s\n", key.(string), b)
		return true
	}
	sub.km.key.Range(dumper)
	log.Printf("\trestarted %t\n", sub.km.restarted)
}

func (sub *Subscription) Get(key string) (interface{}, error) {
	m, ok := sub.km.key.Load(key)
	if ok {
		return m, nil
	} else {
		name := sub.nameString()
		errStr := fmt.Sprintf("Get(%s) unknown key %s", name, key)
		return nil, errors.New(errStr)
	}
}

// Enumerate all the key, value for the collection
func (sub *Subscription) GetAll() map[string]interface{} {
	result := make(map[string]interface{})
	assigner := func(key, val interface{}) bool {
		result[key.(string)] = val
		return true
	}
	sub.km.key.Range(assigner)
	return result
}

func (sub *Subscription) Restarted() bool {
	return sub.km.restarted
}
