// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.0
// source: key.proto

package tokens

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

type Key struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"` // Hashed value.
	Salt []byte `protobuf:"bytes,2,opt,name=salt,proto3" json:"salt,omitempty"` // Additional data hashed with value.
	// Types that are assignable to State:
	//
	//	*Key_Argon2
	State isKey_State `protobuf_oneof:"state"`
}

func (x *Key) Reset() {
	*x = Key{}
	if protoimpl.UnsafeEnabled {
		mi := &file_key_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Key) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Key) ProtoMessage() {}

func (x *Key) ProtoReflect() protoreflect.Message {
	mi := &file_key_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Key.ProtoReflect.Descriptor instead.
func (*Key) Descriptor() ([]byte, []int) {
	return file_key_proto_rawDescGZIP(), []int{0}
}

func (x *Key) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Key) GetSalt() []byte {
	if x != nil {
		return x.Salt
	}
	return nil
}

func (m *Key) GetState() isKey_State {
	if m != nil {
		return m.State
	}
	return nil
}

func (x *Key) GetArgon2() *Argon2State {
	if x, ok := x.GetState().(*Key_Argon2); ok {
		return x.Argon2
	}
	return nil
}

type isKey_State interface {
	isKey_State()
}

type Key_Argon2 struct {
	Argon2 *Argon2State `protobuf:"bytes,15,opt,name=argon2,proto3,oneof"` // Argon2
}

func (*Key_Argon2) isKey_State() {}

type Argon2State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parallelism uint32 `protobuf:"varint,1,opt,name=parallelism,proto3" json:"parallelism,omitempty"`                 // Degree of parallelism (i.e. number of threads).
	TagLength   uint32 `protobuf:"varint,2,opt,name=tag_length,json=tagLength,proto3" json:"tag_length,omitempty"`    // Desired number of returned bytes.
	MemorySize  uint32 `protobuf:"varint,3,opt,name=memory_size,json=memorySize,proto3" json:"memory_size,omitempty"` // Amount of memory (in kibibytes) to use.
	Iterations  uint32 `protobuf:"varint,4,opt,name=iterations,proto3" json:"iterations,omitempty"`                   // Number of iterations to perform.
	Version     uint32 `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`                         // The current version is 0x13 (19 decimal).
	// 0: Argon2d
	// 1: Argon2i
	// 2: Argon2id
	HashType uint32 `protobuf:"varint,6,opt,name=hash_type,json=hashType,proto3" json:"hash_type,omitempty"`
}

func (x *Argon2State) Reset() {
	*x = Argon2State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_key_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Argon2State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Argon2State) ProtoMessage() {}

func (x *Argon2State) ProtoReflect() protoreflect.Message {
	mi := &file_key_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Argon2State.ProtoReflect.Descriptor instead.
func (*Argon2State) Descriptor() ([]byte, []int) {
	return file_key_proto_rawDescGZIP(), []int{1}
}

func (x *Argon2State) GetParallelism() uint32 {
	if x != nil {
		return x.Parallelism
	}
	return 0
}

func (x *Argon2State) GetTagLength() uint32 {
	if x != nil {
		return x.TagLength
	}
	return 0
}

func (x *Argon2State) GetMemorySize() uint32 {
	if x != nil {
		return x.MemorySize
	}
	return 0
}

func (x *Argon2State) GetIterations() uint32 {
	if x != nil {
		return x.Iterations
	}
	return 0
}

func (x *Argon2State) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Argon2State) GetHashType() uint32 {
	if x != nil {
		return x.HashType
	}
	return 0
}

var File_key_proto protoreflect.FileDescriptor

var file_key_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x22, 0x65, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x61, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x61,
	0x6c, 0x74, 0x12, 0x2d, 0x0a, 0x06, 0x61, 0x72, 0x67, 0x6f, 0x6e, 0x32, 0x18, 0x0f, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x2e, 0x41, 0x72, 0x67, 0x6f,
	0x6e, 0x32, 0x53, 0x74, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x06, 0x61, 0x72, 0x67, 0x6f, 0x6e,
	0x32, 0x42, 0x07, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xc6, 0x01, 0x0a, 0x0b, 0x41,
	0x72, 0x67, 0x6f, 0x6e, 0x32, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61,
	0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x69, 0x73, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0b, 0x70, 0x61, 0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x69, 0x73, 0x6d, 0x12, 0x1d, 0x0a, 0x0a,
	0x74, 0x61, 0x67, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x74, 0x61, 0x67, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x6d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a,
	0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x68, 0x61, 0x73, 0x68, 0x54,
	0x79, 0x70, 0x65, 0x42, 0x19, 0x5a, 0x17, 0x6b, 0x68, 0x65, 0x70, 0x72, 0x69, 0x2e, 0x64, 0x65,
	0x76, 0x2f, 0x68, 0x6f, 0x72, 0x75, 0x73, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_key_proto_rawDescOnce sync.Once
	file_key_proto_rawDescData = file_key_proto_rawDesc
)

func file_key_proto_rawDescGZIP() []byte {
	file_key_proto_rawDescOnce.Do(func() {
		file_key_proto_rawDescData = protoimpl.X.CompressGZIP(file_key_proto_rawDescData)
	})
	return file_key_proto_rawDescData
}

var file_key_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_key_proto_goTypes = []interface{}{
	(*Key)(nil),         // 0: tokens.Key
	(*Argon2State)(nil), // 1: tokens.Argon2State
}
var file_key_proto_depIdxs = []int32{
	1, // 0: tokens.Key.argon2:type_name -> tokens.Argon2State
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_key_proto_init() }
func file_key_proto_init() {
	if File_key_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_key_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Key); i {
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
		file_key_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Argon2State); i {
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
	file_key_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Key_Argon2)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_key_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_key_proto_goTypes,
		DependencyIndexes: file_key_proto_depIdxs,
		MessageInfos:      file_key_proto_msgTypes,
	}.Build()
	File_key_proto = out.File
	file_key_proto_rawDesc = nil
	file_key_proto_goTypes = nil
	file_key_proto_depIdxs = nil
}
