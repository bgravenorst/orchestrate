// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: pkg/types/ethereum/transaction.proto

package ethereum

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Transaction
type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	// QUANTITY - Integer of a nonce.
	Nonce string `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	// DATA (20 Bytes) - The address of the receiver. null when it’s a contract creation transaction.
	// e.g. 0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C
	To string `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	// QUANTITY - Integer of the value sent with this transaction.
	// e.g 0xaf23
	Value string `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	// QUANTITY - Integer of the gas provided for the transaction execution.
	Gas string `protobuf:"bytes,5,opt,name=gas,proto3" json:"gas,omitempty"`
	// QUANTITY - Integer of the gas price used for each paid gas.
	// e.g 0xaf23b
	GasPrice string `protobuf:"bytes,6,opt,name=gas_price,json=gasPrice,proto3" json:"gas_price,omitempty"`
	// DATA - Hash of the method signature (4 bytes) followed by encoded parameters.
	// e.g 0xa9059cbb000000000000000000000000ff778b716fc07d98839f48ddb88d8be583beb684000000000000000000000000000000000000000000000000002386f26fc10000
	Data string `protobuf:"bytes,7,opt,name=data,proto3" json:"data,omitempty"`
	// DATA - The signed, RLP encoded transaction
	Raw string `protobuf:"bytes,8,opt,name=raw,proto3" json:"raw,omitempty"`
	// DATA (32 Bytes) - Hash of the transaction.
	// e.g. 0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210
	TxHash string `protobuf:"bytes,9,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_types_ethereum_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_types_ethereum_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_pkg_types_ethereum_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *Transaction) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *Transaction) GetNonce() string {
	if x != nil {
		return x.Nonce
	}
	return ""
}

func (x *Transaction) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *Transaction) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Transaction) GetGas() string {
	if x != nil {
		return x.Gas
	}
	return ""
}

func (x *Transaction) GetGasPrice() string {
	if x != nil {
		return x.GasPrice
	}
	return ""
}

func (x *Transaction) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *Transaction) GetRaw() string {
	if x != nil {
		return x.Raw
	}
	return ""
}

func (x *Transaction) GetTxHash() string {
	if x != nil {
		return x.TxHash
	}
	return ""
}

var File_pkg_types_ethereum_transaction_proto protoreflect.FileDescriptor

var file_pkg_types_ethereum_transaction_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x65, 0x74, 0x68, 0x65,
	0x72, 0x65, 0x75, 0x6d, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d,
	0x22, 0xcb, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x66, 0x72, 0x6f, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x6f,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x67, 0x61, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x67,
	0x61, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x61, 0x73, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x61, 0x73, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x61, 0x77, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x72, 0x61, 0x77, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x78, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x78, 0x48, 0x61, 0x73, 0x68, 0x42, 0x35,
	0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6f, 0x6e,
	0x73, 0x65, 0x6e, 0x53, 0x79, 0x73, 0x2f, 0x6f, 0x72, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x65, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_types_ethereum_transaction_proto_rawDescOnce sync.Once
	file_pkg_types_ethereum_transaction_proto_rawDescData = file_pkg_types_ethereum_transaction_proto_rawDesc
)

func file_pkg_types_ethereum_transaction_proto_rawDescGZIP() []byte {
	file_pkg_types_ethereum_transaction_proto_rawDescOnce.Do(func() {
		file_pkg_types_ethereum_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_types_ethereum_transaction_proto_rawDescData)
	})
	return file_pkg_types_ethereum_transaction_proto_rawDescData
}

var file_pkg_types_ethereum_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_types_ethereum_transaction_proto_goTypes = []interface{}{
	(*Transaction)(nil), // 0: ethereum.Transaction
}
var file_pkg_types_ethereum_transaction_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_types_ethereum_transaction_proto_init() }
func file_pkg_types_ethereum_transaction_proto_init() {
	if File_pkg_types_ethereum_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_types_ethereum_transaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_types_ethereum_transaction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_types_ethereum_transaction_proto_goTypes,
		DependencyIndexes: file_pkg_types_ethereum_transaction_proto_depIdxs,
		MessageInfos:      file_pkg_types_ethereum_transaction_proto_msgTypes,
	}.Build()
	File_pkg_types_ethereum_transaction_proto = out.File
	file_pkg_types_ethereum_transaction_proto_rawDesc = nil
	file_pkg_types_ethereum_transaction_proto_goTypes = nil
	file_pkg_types_ethereum_transaction_proto_depIdxs = nil
}