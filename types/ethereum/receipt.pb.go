// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types/ethereum/receipt.proto

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

// Ethereum Log
type Log struct {
	// DATA (20 Bytes) - Address from which log originated
	// e.g 0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C
	Address *Account `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Array of DATA (32 Bytes) - Array of 0 to 4 indexed log arguments
	// e.g. 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	Topics []*Hash `protobuf:"bytes,2,rep,name=topics,proto3" json:"topics,omitempty"`
	// DATA - Non-indexed arguments of the log
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	// Log decoded data
	Event       string            `protobuf:"bytes,4,opt,name=event,proto3" json:"event,omitempty"`
	DecodedData map[string]string `protobuf:"bytes,5,rep,name=decoded_data,json=decodedData,proto3" json:"decoded_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// QUANTITY - Block number where this transaction was in
	BlockNumber uint64 `protobuf:"varint,6,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// DATA (32 Bytes) - Hash of the transaction.
	// e.g. 0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e
	TxHash *Hash `protobuf:"bytes,7,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	// QUANTITY - Integer of the transactions index position in the block.
	TxIndex uint64 `protobuf:"varint,8,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// DATA (32 Bytes) - Hash of the block where this transaction was in.
	// e.g. 0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056
	BlockHash *Hash `protobuf:"bytes,9,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// QUANTITY - Integer of the log index position in the block
	Index uint64 `protobuf:"varint,10,opt,name=index,proto3" json:"index,omitempty"`
	// Removed field is true if this log was reverted due to a chain reorganisation.
	Removed              bool     `protobuf:"varint,11,opt,name=removed,proto3" json:"removed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}
func (*Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_797a9f8629b2cab8, []int{0}
}

func (m *Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Log.Unmarshal(m, b)
}
func (m *Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Log.Marshal(b, m, deterministic)
}
func (m *Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Log.Merge(m, src)
}
func (m *Log) XXX_Size() int {
	return xxx_messageInfo_Log.Size(m)
}
func (m *Log) XXX_DiscardUnknown() {
	xxx_messageInfo_Log.DiscardUnknown(m)
}

var xxx_messageInfo_Log proto.InternalMessageInfo

func (m *Log) GetAddress() *Account {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Log) GetTopics() []*Hash {
	if m != nil {
		return m.Topics
	}
	return nil
}

func (m *Log) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Log) GetEvent() string {
	if m != nil {
		return m.Event
	}
	return ""
}

func (m *Log) GetDecodedData() map[string]string {
	if m != nil {
		return m.DecodedData
	}
	return nil
}

func (m *Log) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *Log) GetTxHash() *Hash {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *Log) GetTxIndex() uint64 {
	if m != nil {
		return m.TxIndex
	}
	return 0
}

func (m *Log) GetBlockHash() *Hash {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *Log) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Log) GetRemoved() bool {
	if m != nil {
		return m.Removed
	}
	return false
}

// Transaction Receipt
type Receipt struct {
	// HASH (32 Bytes) - Hash of the transaction.
	// e.g 0x3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e
	TxHash *Hash `protobuf:"bytes,1,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	// HASH (32 Bytes) - Hash of the block where this transaction was in.
	// e.g. 0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056
	BlockHash *Hash `protobuf:"bytes,2,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// QUANTITY - Block number where this transaction was in.
	BlockNumber uint64 `protobuf:"varint,3,opt,name=block_number,json=blockNumber,proto3" json:"block_number,omitempty"`
	// QUANTITY - Integer of the transactions index position in the block.
	TxIndex uint64 `protobuf:"varint,4,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// DATA (20 Bytes) - The contract address created, if the transaction was a contract creation
	// e.g 0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C
	ContractAddress *Account `protobuf:"bytes,6,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	// DATA (32 Bytes) - State root hash after executing transaction
	// e.g. 0x656c34545f90a730a19008c0e7a7cd4fb3895064b48d6d69761bd5abad681056
	PostState []byte `protobuf:"bytes,7,opt,name=post_state,json=postState,proto3" json:"post_state,omitempty"`
	// QUANTITY - 0 indicates transaction failure , 1 indicates transaction success.
	Status uint64 `protobuf:"varint,8,opt,name=status,proto3" json:"status,omitempty"`
	// DATA (256 Bytes) - Bloom filter of logs/events generated by contracts during transaction execution.
	Bloom []byte `protobuf:"bytes,10,opt,name=bloom,proto3" json:"bloom,omitempty"`
	// Array - Array of log objects, which this transaction generated.
	Logs []*Log `protobuf:"bytes,11,rep,name=logs,proto3" json:"logs,omitempty"`
	// QUANTITY - The amount of gas used by this specific transaction alone.
	GasUsed uint64 `protobuf:"varint,12,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// QUANTITY - The total amount of gas used when this transaction was executed in the block.
	CumulativeGasUsed    uint64   `protobuf:"varint,13,opt,name=cumulative_gas_used,json=cumulativeGasUsed,proto3" json:"cumulative_gas_used,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Receipt) Reset()         { *m = Receipt{} }
func (m *Receipt) String() string { return proto.CompactTextString(m) }
func (*Receipt) ProtoMessage()    {}
func (*Receipt) Descriptor() ([]byte, []int) {
	return fileDescriptor_797a9f8629b2cab8, []int{1}
}

func (m *Receipt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receipt.Unmarshal(m, b)
}
func (m *Receipt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receipt.Marshal(b, m, deterministic)
}
func (m *Receipt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receipt.Merge(m, src)
}
func (m *Receipt) XXX_Size() int {
	return xxx_messageInfo_Receipt.Size(m)
}
func (m *Receipt) XXX_DiscardUnknown() {
	xxx_messageInfo_Receipt.DiscardUnknown(m)
}

var xxx_messageInfo_Receipt proto.InternalMessageInfo

func (m *Receipt) GetTxHash() *Hash {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *Receipt) GetBlockHash() *Hash {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *Receipt) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *Receipt) GetTxIndex() uint64 {
	if m != nil {
		return m.TxIndex
	}
	return 0
}

func (m *Receipt) GetContractAddress() *Account {
	if m != nil {
		return m.ContractAddress
	}
	return nil
}

func (m *Receipt) GetPostState() []byte {
	if m != nil {
		return m.PostState
	}
	return nil
}

func (m *Receipt) GetStatus() uint64 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Receipt) GetBloom() []byte {
	if m != nil {
		return m.Bloom
	}
	return nil
}

func (m *Receipt) GetLogs() []*Log {
	if m != nil {
		return m.Logs
	}
	return nil
}

func (m *Receipt) GetGasUsed() uint64 {
	if m != nil {
		return m.GasUsed
	}
	return 0
}

func (m *Receipt) GetCumulativeGasUsed() uint64 {
	if m != nil {
		return m.CumulativeGasUsed
	}
	return 0
}

func init() {
	proto.RegisterType((*Log)(nil), "ethereum.Log")
	proto.RegisterMapType((map[string]string)(nil), "ethereum.Log.DecodedDataEntry")
	proto.RegisterType((*Receipt)(nil), "ethereum.Receipt")
}

func init() { proto.RegisterFile("types/ethereum/receipt.proto", fileDescriptor_797a9f8629b2cab8) }

var fileDescriptor_797a9f8629b2cab8 = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0x95, 0xe3, 0x34, 0x4e, 0xc6, 0xee, 0xef, 0xd7, 0x2e, 0x15, 0xda, 0x56, 0x80, 0xdc, 0x1e,
	0xc0, 0x12, 0xaa, 0x2d, 0x95, 0x0b, 0x42, 0x08, 0xa9, 0x50, 0xa0, 0x48, 0x15, 0x87, 0xad, 0xb8,
	0x70, 0xb1, 0xd6, 0xeb, 0xc1, 0xb1, 0x62, 0x7b, 0x23, 0xef, 0x3a, 0x4a, 0xbe, 0x09, 0x5f, 0x8a,
	0xef, 0x84, 0xbc, 0x8e, 0x49, 0x93, 0xf2, 0xe7, 0x94, 0x79, 0x33, 0x2f, 0x2f, 0x79, 0x33, 0x6f,
	0xe1, 0x91, 0x5e, 0xcd, 0x51, 0x45, 0xa8, 0xa7, 0x58, 0x63, 0x53, 0x46, 0x35, 0x0a, 0xcc, 0xe7,
	0x3a, 0x9c, 0xd7, 0x52, 0x4b, 0x32, 0xee, 0xfb, 0x27, 0xc7, 0x3b, 0xbc, 0x84, 0x2b, 0xec, 0x48,
	0x67, 0x3f, 0x6c, 0xb0, 0x6f, 0x64, 0x46, 0x9e, 0x83, 0xc3, 0xd3, 0xb4, 0x46, 0xa5, 0xa8, 0xe5,
	0x5b, 0x81, 0x7b, 0x71, 0x18, 0xf6, 0xf4, 0xf0, 0x52, 0x08, 0xd9, 0x54, 0x9a, 0xf5, 0x0c, 0xf2,
	0x14, 0x46, 0x5a, 0xce, 0x73, 0xa1, 0xe8, 0xc0, 0xb7, 0x03, 0xf7, 0xe2, 0xbf, 0x0d, 0xf7, 0x9a,
	0xab, 0x29, 0x5b, 0x4f, 0x09, 0x81, 0x61, 0xca, 0x35, 0xa7, 0xb6, 0x6f, 0x05, 0x1e, 0x33, 0x35,
	0x39, 0x82, 0x3d, 0x5c, 0x60, 0xa5, 0xe9, 0xd0, 0xb7, 0x82, 0x09, 0xeb, 0x00, 0xb9, 0x04, 0x2f,
	0x45, 0x21, 0x53, 0x4c, 0x63, 0xf3, 0x8d, 0x3d, 0xa3, 0xfb, 0x64, 0xa3, 0x7b, 0x23, 0xb3, 0xf0,
	0xaa, 0x63, 0x5c, 0x71, 0xcd, 0xdf, 0x57, 0xba, 0x5e, 0x31, 0x37, 0xdd, 0x74, 0xc8, 0x29, 0x78,
	0x49, 0x21, 0xc5, 0x2c, 0xae, 0x9a, 0x32, 0xc1, 0x9a, 0x8e, 0x7c, 0x2b, 0x18, 0x32, 0xd7, 0xf4,
	0x3e, 0x9b, 0x16, 0x79, 0x06, 0x8e, 0x5e, 0xc6, 0x53, 0xae, 0xa6, 0xd4, 0x31, 0x26, 0xef, 0xff,
	0xf1, 0x65, 0xfb, 0x49, 0x8e, 0x61, 0xac, 0x97, 0x71, 0x5e, 0xa5, 0xb8, 0xa4, 0x63, 0xa3, 0xe3,
	0xe8, 0xe5, 0xa7, 0x16, 0x92, 0x73, 0x80, 0xee, 0x67, 0x8c, 0xcc, 0xe4, 0xb7, 0x32, 0x13, 0xc3,
	0x30, 0x4a, 0x47, 0xb0, 0xd7, 0xc9, 0x80, 0x91, 0xe9, 0x00, 0xa1, 0xe0, 0xd4, 0x58, 0xca, 0x05,
	0xa6, 0xd4, 0xf5, 0xad, 0x60, 0xcc, 0x7a, 0x78, 0xf2, 0x06, 0x0e, 0x76, 0x6d, 0x92, 0x03, 0xb0,
	0x67, 0xb8, 0x32, 0x77, 0x99, 0xb0, 0xb6, 0x6c, 0x55, 0x17, 0xbc, 0x68, 0x90, 0x0e, 0xba, 0x25,
	0x1a, 0xf0, 0x6a, 0xf0, 0xd2, 0x3a, 0xfb, 0x6e, 0x83, 0xc3, 0xba, 0x18, 0xdc, 0xb5, 0x6b, 0xfd,
	0xd5, 0xee, 0xb6, 0xa7, 0xc1, 0xbf, 0x3c, 0xed, 0x6e, 0xda, 0xbe, 0xbf, 0xe9, 0xbb, 0x0b, 0x1c,
	0x6e, 0x2f, 0xf0, 0x35, 0x1c, 0x08, 0x59, 0xe9, 0x9a, 0x0b, 0x1d, 0xf7, 0x91, 0x1b, 0xfd, 0x29,
	0x72, 0xff, 0xf7, 0xd4, 0xcb, 0x75, 0xf4, 0x1e, 0x03, 0xcc, 0xa5, 0xd2, 0xb1, 0xd2, 0x5c, 0xa3,
	0xb9, 0xa2, 0xc7, 0x26, 0x6d, 0xe7, 0xb6, 0x6d, 0x90, 0x87, 0x30, 0x6a, 0x27, 0x8d, 0x5a, 0x9f,
	0x6d, 0x8d, 0xda, 0x85, 0x25, 0x85, 0x94, 0xa5, 0x39, 0x83, 0xc7, 0x3a, 0x40, 0x4e, 0x61, 0x58,
	0xc8, 0x4c, 0x51, 0xd7, 0xa4, 0x6d, 0x7f, 0x2b, 0x6d, 0xcc, 0x8c, 0x5a, 0x23, 0x19, 0x57, 0x71,
	0xa3, 0x30, 0xa5, 0x5e, 0x67, 0x24, 0xe3, 0xea, 0x8b, 0xc2, 0x94, 0x84, 0xf0, 0x40, 0x34, 0x65,
	0x53, 0x70, 0x9d, 0x2f, 0x30, 0xfe, 0xc5, 0xda, 0x37, 0xac, 0xc3, 0xcd, 0xe8, 0x63, 0xc7, 0x7f,
	0x7b, 0xfd, 0xf5, 0x43, 0x96, 0xeb, 0x82, 0x27, 0xa1, 0x90, 0x65, 0xf4, 0x4e, 0x56, 0x0a, 0xab,
	0xdb, 0x95, 0x8a, 0x44, 0x91, 0x63, 0xa5, 0xa3, 0x6f, 0x75, 0x24, 0x64, 0x8d, 0xe7, 0x4a, 0x73,
	0x31, 0x33, 0xa5, 0xa9, 0xc2, 0x2c, 0xd7, 0xd1, 0xf6, 0xfb, 0x4d, 0x46, 0xe6, 0xed, 0xbe, 0xf8,
	0x19, 0x00, 0x00, 0xff, 0xff, 0xe8, 0x53, 0x9b, 0xba, 0x00, 0x04, 0x00, 0x00,
}