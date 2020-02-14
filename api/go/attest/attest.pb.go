// Code generated by protoc-gen-go. DO NOT EDIT.
// source: attest.proto

package attest

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ZAttestReqType int32

const (
	ZAttestReqType_ATTEST_REQ_CERT  ZAttestReqType = 0
	ZAttestReqType_ATTEST_REQ_NONCE ZAttestReqType = 1
	ZAttestReqType_ATTEST_REQ_QUOTE ZAttestReqType = 2
)

var ZAttestReqType_name = map[int32]string{
	0: "ATTEST_REQ_CERT",
	1: "ATTEST_REQ_NONCE",
	2: "ATTEST_REQ_QUOTE",
}

var ZAttestReqType_value = map[string]int32{
	"ATTEST_REQ_CERT":  0,
	"ATTEST_REQ_NONCE": 1,
	"ATTEST_REQ_QUOTE": 2,
}

func (x ZAttestReqType) String() string {
	return proto.EnumName(ZAttestReqType_name, int32(x))
}

func (ZAttestReqType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{0}
}

type ZAttestRespType int32

const (
	ZAttestRespType_ATTEST_RESP_CERT       ZAttestRespType = 0
	ZAttestRespType_ATTEST_RESP_NONCE      ZAttestRespType = 1
	ZAttestRespType_ATTEST_RESP_QUOTE_RESP ZAttestRespType = 2
)

var ZAttestRespType_name = map[int32]string{
	0: "ATTEST_RESP_CERT",
	1: "ATTEST_RESP_NONCE",
	2: "ATTEST_RESP_QUOTE_RESP",
}

var ZAttestRespType_value = map[string]int32{
	"ATTEST_RESP_CERT":       0,
	"ATTEST_RESP_NONCE":      1,
	"ATTEST_RESP_QUOTE_RESP": 2,
}

func (x ZAttestRespType) String() string {
	return proto.EnumName(ZAttestRespType_name, int32(x))
}

func (ZAttestRespType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{1}
}

type ZAttestResponseCode int32

const (
	ZAttestResponseCode_ATTEST_RESPONSE_SUCCESS ZAttestResponseCode = 0
	ZAttestResponseCode_ATTEST_RESPONSE_FAILURE ZAttestResponseCode = 1
)

var ZAttestResponseCode_name = map[int32]string{
	0: "ATTEST_RESPONSE_SUCCESS",
	1: "ATTEST_RESPONSE_FAILURE",
}

var ZAttestResponseCode_value = map[string]int32{
	"ATTEST_RESPONSE_SUCCESS": 0,
	"ATTEST_RESPONSE_FAILURE": 1,
}

func (x ZAttestResponseCode) String() string {
	return proto.EnumName(ZAttestResponseCode_name, int32(x))
}

func (ZAttestResponseCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{2}
}

type ZEveCertHashType int32

const (
	ZEveCertHashType_HASH_NONE           ZEveCertHashType = 0
	ZEveCertHashType_HASH_SHA256_16bytes ZEveCertHashType = 1
)

var ZEveCertHashType_name = map[int32]string{
	0: "HASH_NONE",
	1: "HASH_SHA256_16bytes",
}

var ZEveCertHashType_value = map[string]int32{
	"HASH_NONE":           0,
	"HASH_SHA256_16bytes": 1,
}

func (x ZEveCertHashType) String() string {
	return proto.EnumName(ZEveCertHashType_name, int32(x))
}

func (ZEveCertHashType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{3}
}

type ZEveCertType int32

const (
	ZEveCertType_CERT_TYPE_DEVICE_ONBOARDING         ZEveCertType = 0
	ZEveCertType_CERT_TYPE_DEVICE_RESTRICTED_SIGNING ZEveCertType = 1
	ZEveCertType_CERT_TYPE_DEVICE_ENDORSEMENT_RSA    ZEveCertType = 2
	ZEveCertType_CERT_TYPE_DEVICE_ECDH_EXCHANGE      ZEveCertType = 3
)

var ZEveCertType_name = map[int32]string{
	0: "CERT_TYPE_DEVICE_ONBOARDING",
	1: "CERT_TYPE_DEVICE_RESTRICTED_SIGNING",
	2: "CERT_TYPE_DEVICE_ENDORSEMENT_RSA",
	3: "CERT_TYPE_DEVICE_ECDH_EXCHANGE",
}

var ZEveCertType_value = map[string]int32{
	"CERT_TYPE_DEVICE_ONBOARDING":         0,
	"CERT_TYPE_DEVICE_RESTRICTED_SIGNING": 1,
	"CERT_TYPE_DEVICE_ENDORSEMENT_RSA":    2,
	"CERT_TYPE_DEVICE_ECDH_EXCHANGE":      3,
}

func (x ZEveCertType) String() string {
	return proto.EnumName(ZEveCertType_name, int32(x))
}

func (ZEveCertType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{4}
}

// This is the request payload for POST /api/v2/edgeDevice/id/<uuid>/attest
// The message is assumed to be protected by signing envelope
type ZAttestReq struct {
	ReqType              ZAttestReqType `protobuf:"varint,1,opt,name=reqType,proto3,enum=ZAttestReqType" json:"reqType,omitempty"`
	Quote                *ZAttestQuote  `protobuf:"bytes,2,opt,name=quote,proto3" json:"quote,omitempty"`
	Certs                []*ZEveCert    `protobuf:"bytes,3,rep,name=certs,proto3" json:"certs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ZAttestReq) Reset()         { *m = ZAttestReq{} }
func (m *ZAttestReq) String() string { return proto.CompactTextString(m) }
func (*ZAttestReq) ProtoMessage()    {}
func (*ZAttestReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{0}
}

func (m *ZAttestReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZAttestReq.Unmarshal(m, b)
}
func (m *ZAttestReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZAttestReq.Marshal(b, m, deterministic)
}
func (m *ZAttestReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZAttestReq.Merge(m, src)
}
func (m *ZAttestReq) XXX_Size() int {
	return xxx_messageInfo_ZAttestReq.Size(m)
}
func (m *ZAttestReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ZAttestReq.DiscardUnknown(m)
}

var xxx_messageInfo_ZAttestReq proto.InternalMessageInfo

func (m *ZAttestReq) GetReqType() ZAttestReqType {
	if m != nil {
		return m.ReqType
	}
	return ZAttestReqType_ATTEST_REQ_CERT
}

func (m *ZAttestReq) GetQuote() *ZAttestQuote {
	if m != nil {
		return m.Quote
	}
	return nil
}

func (m *ZAttestReq) GetCerts() []*ZEveCert {
	if m != nil {
		return m.Certs
	}
	return nil
}

// This is the response payload for POST /api/v2/edgeDevice/id/<uuid>/attest
// The message is assumed to be protected by signing envelope
type ZAttestResponse struct {
	RespType             ZAttestRespType   `protobuf:"varint,1,opt,name=respType,proto3,enum=ZAttestRespType" json:"respType,omitempty"`
	Nonce                *ZAttestNonceResp `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	QuoteResp            *ZAttestQuoteResp `protobuf:"bytes,3,opt,name=quoteResp,proto3" json:"quoteResp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ZAttestResponse) Reset()         { *m = ZAttestResponse{} }
func (m *ZAttestResponse) String() string { return proto.CompactTextString(m) }
func (*ZAttestResponse) ProtoMessage()    {}
func (*ZAttestResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{1}
}

func (m *ZAttestResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZAttestResponse.Unmarshal(m, b)
}
func (m *ZAttestResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZAttestResponse.Marshal(b, m, deterministic)
}
func (m *ZAttestResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZAttestResponse.Merge(m, src)
}
func (m *ZAttestResponse) XXX_Size() int {
	return xxx_messageInfo_ZAttestResponse.Size(m)
}
func (m *ZAttestResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ZAttestResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ZAttestResponse proto.InternalMessageInfo

func (m *ZAttestResponse) GetRespType() ZAttestRespType {
	if m != nil {
		return m.RespType
	}
	return ZAttestRespType_ATTEST_RESP_CERT
}

func (m *ZAttestResponse) GetNonce() *ZAttestNonceResp {
	if m != nil {
		return m.Nonce
	}
	return nil
}

func (m *ZAttestResponse) GetQuoteResp() *ZAttestQuoteResp {
	if m != nil {
		return m.QuoteResp
	}
	return nil
}

type ZAttestNonceResp struct {
	Nonce                []byte   `protobuf:"bytes,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ZAttestNonceResp) Reset()         { *m = ZAttestNonceResp{} }
func (m *ZAttestNonceResp) String() string { return proto.CompactTextString(m) }
func (*ZAttestNonceResp) ProtoMessage()    {}
func (*ZAttestNonceResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{2}
}

func (m *ZAttestNonceResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZAttestNonceResp.Unmarshal(m, b)
}
func (m *ZAttestNonceResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZAttestNonceResp.Marshal(b, m, deterministic)
}
func (m *ZAttestNonceResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZAttestNonceResp.Merge(m, src)
}
func (m *ZAttestNonceResp) XXX_Size() int {
	return xxx_messageInfo_ZAttestNonceResp.Size(m)
}
func (m *ZAttestNonceResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ZAttestNonceResp.DiscardUnknown(m)
}

var xxx_messageInfo_ZAttestNonceResp proto.InternalMessageInfo

func (m *ZAttestNonceResp) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

type ZAttestQuote struct {
	AttestData []byte `protobuf:"bytes,1,opt,name=attestData,proto3" json:"attestData,omitempty"`
	//nonce is included in attestData
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ZAttestQuote) Reset()         { *m = ZAttestQuote{} }
func (m *ZAttestQuote) String() string { return proto.CompactTextString(m) }
func (*ZAttestQuote) ProtoMessage()    {}
func (*ZAttestQuote) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{3}
}

func (m *ZAttestQuote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZAttestQuote.Unmarshal(m, b)
}
func (m *ZAttestQuote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZAttestQuote.Marshal(b, m, deterministic)
}
func (m *ZAttestQuote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZAttestQuote.Merge(m, src)
}
func (m *ZAttestQuote) XXX_Size() int {
	return xxx_messageInfo_ZAttestQuote.Size(m)
}
func (m *ZAttestQuote) XXX_DiscardUnknown() {
	xxx_messageInfo_ZAttestQuote.DiscardUnknown(m)
}

var xxx_messageInfo_ZAttestQuote proto.InternalMessageInfo

func (m *ZAttestQuote) GetAttestData() []byte {
	if m != nil {
		return m.AttestData
	}
	return nil
}

func (m *ZAttestQuote) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type ZAttestQuoteResp struct {
	Response             ZAttestResponseCode `protobuf:"varint,1,opt,name=response,proto3,enum=ZAttestResponseCode" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ZAttestQuoteResp) Reset()         { *m = ZAttestQuoteResp{} }
func (m *ZAttestQuoteResp) String() string { return proto.CompactTextString(m) }
func (*ZAttestQuoteResp) ProtoMessage()    {}
func (*ZAttestQuoteResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{4}
}

func (m *ZAttestQuoteResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZAttestQuoteResp.Unmarshal(m, b)
}
func (m *ZAttestQuoteResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZAttestQuoteResp.Marshal(b, m, deterministic)
}
func (m *ZAttestQuoteResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZAttestQuoteResp.Merge(m, src)
}
func (m *ZAttestQuoteResp) XXX_Size() int {
	return xxx_messageInfo_ZAttestQuoteResp.Size(m)
}
func (m *ZAttestQuoteResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ZAttestQuoteResp.DiscardUnknown(m)
}

var xxx_messageInfo_ZAttestQuoteResp proto.InternalMessageInfo

func (m *ZAttestQuoteResp) GetResponse() ZAttestResponseCode {
	if m != nil {
		return m.Response
	}
	return ZAttestResponseCode_ATTEST_RESPONSE_SUCCESS
}

type ZEveCert struct {
	HashAlgo             ZEveCertHashType `protobuf:"varint,1,opt,name=hashAlgo,proto3,enum=ZEveCertHashType" json:"hashAlgo,omitempty"`
	CertHash             []byte           `protobuf:"bytes,2,opt,name=certHash,proto3" json:"certHash,omitempty"`
	Type                 ZEveCertType     `protobuf:"varint,3,opt,name=type,proto3,enum=ZEveCertType" json:"type,omitempty"`
	Cert                 []byte           `protobuf:"bytes,4,opt,name=cert,proto3" json:"cert,omitempty"`
	Attributes           *ZEveCertAttr    `protobuf:"bytes,5,opt,name=attributes,proto3" json:"attributes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ZEveCert) Reset()         { *m = ZEveCert{} }
func (m *ZEveCert) String() string { return proto.CompactTextString(m) }
func (*ZEveCert) ProtoMessage()    {}
func (*ZEveCert) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{5}
}

func (m *ZEveCert) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZEveCert.Unmarshal(m, b)
}
func (m *ZEveCert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZEveCert.Marshal(b, m, deterministic)
}
func (m *ZEveCert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZEveCert.Merge(m, src)
}
func (m *ZEveCert) XXX_Size() int {
	return xxx_messageInfo_ZEveCert.Size(m)
}
func (m *ZEveCert) XXX_DiscardUnknown() {
	xxx_messageInfo_ZEveCert.DiscardUnknown(m)
}

var xxx_messageInfo_ZEveCert proto.InternalMessageInfo

func (m *ZEveCert) GetHashAlgo() ZEveCertHashType {
	if m != nil {
		return m.HashAlgo
	}
	return ZEveCertHashType_HASH_NONE
}

func (m *ZEveCert) GetCertHash() []byte {
	if m != nil {
		return m.CertHash
	}
	return nil
}

func (m *ZEveCert) GetType() ZEveCertType {
	if m != nil {
		return m.Type
	}
	return ZEveCertType_CERT_TYPE_DEVICE_ONBOARDING
}

func (m *ZEveCert) GetCert() []byte {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *ZEveCert) GetAttributes() *ZEveCertAttr {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type ZEveCertAttr struct {
	IsMutable            bool     `protobuf:"varint,1,opt,name=isMutable,proto3" json:"isMutable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ZEveCertAttr) Reset()         { *m = ZEveCertAttr{} }
func (m *ZEveCertAttr) String() string { return proto.CompactTextString(m) }
func (*ZEveCertAttr) ProtoMessage()    {}
func (*ZEveCertAttr) Descriptor() ([]byte, []int) {
	return fileDescriptor_208cf26448842a3f, []int{6}
}

func (m *ZEveCertAttr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZEveCertAttr.Unmarshal(m, b)
}
func (m *ZEveCertAttr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZEveCertAttr.Marshal(b, m, deterministic)
}
func (m *ZEveCertAttr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZEveCertAttr.Merge(m, src)
}
func (m *ZEveCertAttr) XXX_Size() int {
	return xxx_messageInfo_ZEveCertAttr.Size(m)
}
func (m *ZEveCertAttr) XXX_DiscardUnknown() {
	xxx_messageInfo_ZEveCertAttr.DiscardUnknown(m)
}

var xxx_messageInfo_ZEveCertAttr proto.InternalMessageInfo

func (m *ZEveCertAttr) GetIsMutable() bool {
	if m != nil {
		return m.IsMutable
	}
	return false
}

func init() {
	proto.RegisterEnum("ZAttestReqType", ZAttestReqType_name, ZAttestReqType_value)
	proto.RegisterEnum("ZAttestRespType", ZAttestRespType_name, ZAttestRespType_value)
	proto.RegisterEnum("ZAttestResponseCode", ZAttestResponseCode_name, ZAttestResponseCode_value)
	proto.RegisterEnum("ZEveCertHashType", ZEveCertHashType_name, ZEveCertHashType_value)
	proto.RegisterEnum("ZEveCertType", ZEveCertType_name, ZEveCertType_value)
	proto.RegisterType((*ZAttestReq)(nil), "ZAttestReq")
	proto.RegisterType((*ZAttestResponse)(nil), "ZAttestResponse")
	proto.RegisterType((*ZAttestNonceResp)(nil), "ZAttestNonceResp")
	proto.RegisterType((*ZAttestQuote)(nil), "ZAttestQuote")
	proto.RegisterType((*ZAttestQuoteResp)(nil), "ZAttestQuoteResp")
	proto.RegisterType((*ZEveCert)(nil), "ZEveCert")
	proto.RegisterType((*ZEveCertAttr)(nil), "ZEveCertAttr")
}

func init() { proto.RegisterFile("attest.proto", fileDescriptor_208cf26448842a3f) }

var fileDescriptor_208cf26448842a3f = []byte{
	// 672 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x54, 0xd1, 0x52, 0xda, 0x4c,
	0x18, 0x25, 0x20, 0xff, 0x0f, 0x9f, 0xa8, 0xeb, 0x6a, 0x2b, 0xa3, 0x1d, 0xa5, 0xd1, 0x19, 0x29,
	0xa3, 0xa1, 0xa5, 0x53, 0x2f, 0x7a, 0x17, 0xc3, 0x16, 0x98, 0xd1, 0x20, 0x9b, 0xd8, 0x69, 0xb9,
	0xc9, 0x04, 0xd8, 0x02, 0x33, 0x48, 0x30, 0x59, 0x9c, 0xb1, 0x33, 0x7d, 0x90, 0xde, 0xf6, 0x35,
	0xfa, 0x72, 0x9d, 0xdd, 0x24, 0x10, 0xd1, 0xbb, 0xdd, 0x73, 0xce, 0x77, 0xf6, 0x70, 0x76, 0x09,
	0x14, 0x5c, 0xce, 0x59, 0xc0, 0xb5, 0x99, 0xef, 0x71, 0x4f, 0xfd, 0x05, 0xd0, 0xd5, 0x25, 0x40,
	0xd9, 0x3d, 0x7e, 0x07, 0xff, 0xfb, 0xec, 0xde, 0x7e, 0x9c, 0xb1, 0xa2, 0x52, 0x52, 0xca, 0x9b,
	0xb5, 0x2d, 0x6d, 0xc9, 0x0a, 0x98, 0xc6, 0x3c, 0x3e, 0x86, 0xec, 0xfd, 0xdc, 0xe3, 0xac, 0x98,
	0x2e, 0x29, 0xe5, 0xf5, 0xda, 0x46, 0x2c, 0xec, 0x08, 0x90, 0x86, 0x1c, 0x3e, 0x82, 0x6c, 0x9f,
	0xf9, 0x3c, 0x28, 0x66, 0x4a, 0x99, 0xf2, 0x7a, 0x2d, 0xaf, 0x75, 0xc9, 0x03, 0x33, 0x98, 0xcf,
	0x69, 0x88, 0xab, 0xbf, 0x15, 0xd8, 0x5a, 0x9c, 0x10, 0xcc, 0xbc, 0x69, 0xc0, 0xf0, 0x19, 0xe4,
	0x7c, 0x16, 0xcc, 0x12, 0x29, 0x90, 0x96, 0xd0, 0xc8, 0x18, 0x0b, 0x05, 0x3e, 0x85, 0xec, 0xd4,
	0x9b, 0xf6, 0xe3, 0x1c, 0xdb, 0xb1, 0xd4, 0x14, 0xa0, 0xd0, 0xd3, 0x90, 0xc7, 0x55, 0xc8, 0xcb,
	0x50, 0x02, 0x2b, 0x66, 0x9e, 0x8a, 0x3b, 0x31, 0x41, 0x97, 0x1a, 0xb5, 0x0c, 0x68, 0xd5, 0x0b,
	0xef, 0xc6, 0xa7, 0x89, 0x60, 0x85, 0xc8, 0x5a, 0xbd, 0x82, 0x42, 0xd2, 0x08, 0x1f, 0x02, 0x84,
	0x25, 0xd7, 0x5d, 0xee, 0x46, 0xd2, 0x04, 0x82, 0xdf, 0x40, 0x3e, 0x18, 0x0f, 0xa7, 0x2e, 0x9f,
	0xfb, 0x61, 0xee, 0x02, 0x5d, 0x02, 0x6a, 0x7d, 0x71, 0xee, 0x22, 0x16, 0x7e, 0x1f, 0x76, 0x22,
	0xfa, 0x89, 0x3a, 0xd9, 0xd5, 0x56, 0x7a, 0x33, 0xbc, 0x41, 0xd4, 0x8b, 0xd8, 0xa9, 0x7f, 0x15,
	0xc8, 0xc5, 0x6d, 0xe3, 0x73, 0xc8, 0x8d, 0xdc, 0x60, 0xa4, 0x4f, 0x86, 0x5e, 0x34, 0xbe, 0xbd,
	0xb8, 0x8a, 0xa6, 0x1b, 0x8c, 0xc2, 0x4e, 0x63, 0x09, 0xde, 0x87, 0x5c, 0x3f, 0x62, 0xa2, 0x78,
	0x8b, 0x3d, 0x7e, 0x0b, 0x6b, 0x5c, 0xdc, 0x4c, 0x46, 0xda, 0x6c, 0x2c, 0x6c, 0xa4, 0x85, 0xa4,
	0x30, 0x86, 0x35, 0x21, 0x2f, 0xae, 0xc9, 0x51, 0xb9, 0xc6, 0xe7, 0xb2, 0x12, 0x7f, 0xdc, 0x9b,
	0x73, 0x16, 0x14, 0xb3, 0xf1, 0x9b, 0x89, 0x86, 0x75, 0xce, 0x7d, 0x9a, 0x10, 0xa8, 0x67, 0x50,
	0x48, 0x72, 0xa2, 0xb1, 0x71, 0x70, 0x3d, 0xe7, 0x6e, 0x6f, 0x12, 0x16, 0x90, 0xa3, 0x4b, 0xa0,
	0xd2, 0x81, 0xcd, 0xa7, 0xcf, 0x14, 0xef, 0xc0, 0x96, 0x6e, 0xdb, 0xc4, 0xb2, 0x1d, 0x4a, 0x3a,
	0x8e, 0x41, 0xa8, 0x8d, 0x52, 0x78, 0x17, 0x50, 0x02, 0x34, 0xdb, 0xa6, 0x41, 0x90, 0xb2, 0x82,
	0x76, 0x6e, 0xdb, 0x36, 0x41, 0xe9, 0x4a, 0xf7, 0xc9, 0xbb, 0x94, 0x9e, 0x49, 0xa1, 0x75, 0x13,
	0x9b, 0xbe, 0x82, 0xed, 0x24, 0x1a, 0xbb, 0xee, 0xc3, 0xeb, 0x24, 0x2c, 0x6d, 0xe5, 0x12, 0xa5,
	0x2b, 0x6d, 0xd8, 0x79, 0xe1, 0xee, 0xf0, 0x01, 0xec, 0x25, 0x46, 0xda, 0xa6, 0x45, 0x1c, 0xeb,
	0xd6, 0x30, 0x88, 0x65, 0xa1, 0xd4, 0x4b, 0xe4, 0x17, 0xbd, 0x75, 0x75, 0x4b, 0x09, 0x52, 0x2a,
	0x9f, 0x01, 0xad, 0xde, 0x26, 0xde, 0x80, 0x7c, 0x53, 0xb7, 0x9a, 0x22, 0x10, 0x41, 0x29, 0xbc,
	0x07, 0x3b, 0x72, 0x6b, 0x35, 0xf5, 0xda, 0xa7, 0x0b, 0xe7, 0xc3, 0x45, 0xef, 0x91, 0xb3, 0x00,
	0x29, 0x95, 0x3f, 0xca, 0xb2, 0x6a, 0x39, 0x78, 0x04, 0x07, 0xe2, 0xa7, 0x39, 0xf6, 0xf7, 0x1b,
	0xe2, 0xd4, 0xc9, 0xd7, 0x96, 0x41, 0x9c, 0xb6, 0x79, 0xd9, 0xd6, 0x69, 0xbd, 0x65, 0x36, 0x50,
	0x0a, 0x9f, 0xc2, 0xf1, 0x33, 0x01, 0x25, 0x96, 0x4d, 0x5b, 0x86, 0x4d, 0xea, 0x8e, 0xd5, 0x6a,
	0x98, 0x42, 0xa8, 0xe0, 0x13, 0x28, 0x3d, 0x13, 0x12, 0xb3, 0xde, 0xa6, 0x16, 0xb9, 0x26, 0xa6,
	0xed, 0x50, 0x4b, 0x47, 0x69, 0xac, 0xc2, 0xe1, 0x73, 0x95, 0x51, 0x6f, 0x3a, 0xe4, 0x9b, 0xd1,
	0xd4, 0xcd, 0x06, 0x41, 0x99, 0xcb, 0x06, 0x1c, 0xf5, 0xbd, 0x3b, 0xed, 0x27, 0x1b, 0xb0, 0x81,
	0xab, 0xf5, 0x27, 0xde, 0x7c, 0xa0, 0xcd, 0x03, 0xe6, 0x3f, 0x8c, 0xfb, 0x2c, 0xfc, 0x90, 0x75,
	0x4f, 0x86, 0x63, 0x3e, 0x9a, 0xf7, 0xb4, 0xbe, 0x77, 0x57, 0x9d, 0xfc, 0x38, 0x67, 0x83, 0x21,
	0xab, 0xb2, 0x07, 0x56, 0x75, 0x67, 0xe3, 0xea, 0xd0, 0xab, 0x86, 0xff, 0xbe, 0xde, 0x7f, 0x52,
	0xfc, 0xf1, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x85, 0x1f, 0xb3, 0x49, 0x05, 0x05, 0x00, 0x00,
}
