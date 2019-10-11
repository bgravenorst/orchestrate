// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types/ethereum/transaction.proto

package ethereum

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

// Transaction data
type TxData struct {
	// QUANTITY - Integer of a nonce.
	Nonce uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// DATA (20 Bytes) - The address of the receiver. null when it’s a contract creation transaction.
	// e.g. 0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C
	To *Account `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	// QUANTITY - Integer of the value sent with this transaction.
	// e.g 0xaf23
	Value *Quantity `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	// QUANTITY - Integer of the gas provided for the transaction execution.
	Gas uint64 `protobuf:"varint,4,opt,name=gas,proto3" json:"gas,omitempty"`
	// QUANTITY - Integer of the gas price used for each paid gas.
	// e.g 0xaf23b
	GasPrice *Quantity `protobuf:"bytes,5,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	// DATA - Hash of the method signature (4 bytes) followed by encoded parameters.
	// e.g 0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000
	Data                 *Data    `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxData) Reset()         { *m = TxData{} }
func (m *TxData) String() string { return proto.CompactTextString(m) }
func (*TxData) ProtoMessage()    {}
func (*TxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_e8c15aa666d7a06d, []int{0}
}

func (m *TxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxData.Unmarshal(m, b)
}
func (m *TxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxData.Marshal(b, m, deterministic)
}
func (m *TxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxData.Merge(m, src)
}
func (m *TxData) XXX_Size() int {
	return xxx_messageInfo_TxData.Size(m)
}
func (m *TxData) XXX_DiscardUnknown() {
	xxx_messageInfo_TxData.DiscardUnknown(m)
}

var xxx_messageInfo_TxData proto.InternalMessageInfo

func (m *TxData) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *TxData) GetTo() *Account {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *TxData) GetValue() *Quantity {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *TxData) GetGas() uint64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

func (m *TxData) GetGasPrice() *Quantity {
	if m != nil {
		return m.GasPrice
	}
	return nil
}

func (m *TxData) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

// Transaction
type Transaction struct {
	TxData *TxData `protobuf:"bytes,1,opt,name=tx_data,json=txData,proto3" json:"tx_data,omitempty"`
	// DATA - The signed, RLP encoded transaction
	Raw *Data `protobuf:"bytes,2,opt,name=raw,proto3" json:"raw,omitempty"`
	// DATA (32 Bytes) - Hash of the transaction.
	// e.g. 0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210
	Hash                 *Hash    `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_e8c15aa666d7a06d, []int{1}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetTxData() *TxData {
	if m != nil {
		return m.TxData
	}
	return nil
}

func (m *Transaction) GetRaw() *Data {
	if m != nil {
		return m.Raw
	}
	return nil
}

func (m *Transaction) GetHash() *Hash {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*TxData)(nil), "ethereum.TxData")
	proto.RegisterType((*Transaction)(nil), "ethereum.Transaction")
}

func init() { proto.RegisterFile("types/ethereum/transaction.proto", fileDescriptor_e8c15aa666d7a06d) }

var fileDescriptor_e8c15aa666d7a06d = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x95, 0xfe, 0x84, 0xe2, 0x4a, 0xa8, 0x58, 0x0c, 0x86, 0x29, 0x74, 0x2a, 0x03, 0xb1,
	0x04, 0x4f, 0xc0, 0x8f, 0x50, 0x47, 0x30, 0x9d, 0x58, 0xaa, 0x5b, 0x63, 0x12, 0x8b, 0xd4, 0x8e,
	0xec, 0x1b, 0x68, 0x16, 0x1e, 0x92, 0x27, 0x42, 0x71, 0x5a, 0x05, 0x2a, 0xb1, 0x1d, 0xeb, 0x7c,
	0xc7, 0x57, 0xe7, 0x5e, 0x92, 0x60, 0x5d, 0x2a, 0xcf, 0x15, 0xe6, 0xca, 0xa9, 0x6a, 0xcd, 0xd1,
	0x81, 0xf1, 0x20, 0x51, 0x5b, 0x93, 0x96, 0xce, 0xa2, 0xa5, 0xa3, 0x9d, 0x77, 0x76, 0xba, 0xc7,
	0xae, 0xc0, 0xab, 0x16, 0x9a, 0x7e, 0x47, 0x24, 0x5e, 0x6c, 0xee, 0x01, 0x81, 0x9e, 0x90, 0xa1,
	0xb1, 0x46, 0x2a, 0x16, 0x25, 0xd1, 0x6c, 0x20, 0xda, 0x07, 0x3d, 0x27, 0x3d, 0xb4, 0xac, 0x97,
	0x44, 0xb3, 0xf1, 0xd5, 0x71, 0xba, 0xfb, 0x22, 0xbd, 0x91, 0xd2, 0x56, 0x06, 0x45, 0x0f, 0x2d,
	0x9d, 0x91, 0xe1, 0x07, 0x14, 0x95, 0x62, 0xfd, 0x40, 0xd1, 0x8e, 0x7a, 0xaa, 0xc0, 0xa0, 0xc6,
	0x5a, 0xb4, 0x00, 0x9d, 0x90, 0x7e, 0x06, 0x9e, 0x0d, 0xc2, 0x80, 0x46, 0x52, 0x4e, 0x0e, 0x33,
	0xf0, 0xcb, 0xd2, 0x69, 0xa9, 0xd8, 0xf0, 0xdf, 0xfc, 0x28, 0x03, 0xff, 0xd8, 0x30, 0x74, 0x4a,
	0x06, 0xaf, 0x80, 0xc0, 0xe2, 0xc0, 0x1e, 0x75, 0x6c, 0xd3, 0x41, 0x04, 0x6f, 0xfa, 0x45, 0xc6,
	0x8b, 0x6e, 0x1d, 0xf4, 0x82, 0x1c, 0xe0, 0x66, 0x19, 0x52, 0x51, 0x48, 0x4d, 0xba, 0x54, 0xdb,
	0x5d, 0xc4, 0xd8, 0xee, 0x20, 0x21, 0x7d, 0x07, 0x9f, 0xdb, 0xba, 0xfb, 0x9f, 0x37, 0x56, 0x33,
	0x3f, 0x07, 0x9f, 0x6f, 0xbb, 0xfe, 0x42, 0xe6, 0xe0, 0x73, 0x11, 0xbc, 0xdb, 0xf9, 0xcb, 0x43,
	0xa6, 0xb1, 0x80, 0x55, 0x2a, 0xed, 0x9a, 0xdf, 0x59, 0xe3, 0x95, 0x79, 0xae, 0x3d, 0x97, 0x85,
	0x56, 0x06, 0xf9, 0x9b, 0xe3, 0xd2, 0x3a, 0x75, 0xe9, 0x11, 0xe4, 0x7b, 0x90, 0x41, 0xa5, 0x99,
	0x46, 0xfe, 0xf7, 0x52, 0xab, 0x38, 0x5c, 0xe9, 0xfa, 0x27, 0x00, 0x00, 0xff, 0xff, 0x51, 0x12,
	0x8c, 0x96, 0xee, 0x01, 0x00, 0x00,
}