// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: proto/store.proto

package store

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

type PersistedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            []byte `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Stream        []byte `protobuf:"bytes,2,opt,name=stream,proto3" json:"stream,omitempty"`
	Version       uint32 `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Type          string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Data          []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Metadata      []byte `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CausationId   []byte `protobuf:"bytes,7,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId []byte `protobuf:"bytes,8,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	AddedAt       int64  `protobuf:"varint,9,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
}

func (x *PersistedEvent) Reset() {
	*x = PersistedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistedEvent) ProtoMessage() {}

func (x *PersistedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistedEvent.ProtoReflect.Descriptor instead.
func (*PersistedEvent) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{0}
}

func (x *PersistedEvent) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *PersistedEvent) GetStream() []byte {
	if x != nil {
		return x.Stream
	}
	return nil
}

func (x *PersistedEvent) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *PersistedEvent) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *PersistedEvent) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PersistedEvent) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *PersistedEvent) GetCausationId() []byte {
	if x != nil {
		return x.CausationId
	}
	return nil
}

func (x *PersistedEvent) GetCorrelationId() []byte {
	if x != nil {
		return x.CorrelationId
	}
	return nil
}

func (x *PersistedEvent) GetAddedAt() int64 {
	if x != nil {
		return x.AddedAt
	}
	return 0
}

type PersistedStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Events  [][]byte `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
	AddedAt int64    `protobuf:"varint,3,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
}

func (x *PersistedStream) Reset() {
	*x = PersistedStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistedStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistedStream) ProtoMessage() {}

func (x *PersistedStream) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistedStream.ProtoReflect.Descriptor instead.
func (*PersistedStream) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{1}
}

func (x *PersistedStream) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *PersistedStream) GetEvents() [][]byte {
	if x != nil {
		return x.Events
	}
	return nil
}

func (x *PersistedStream) GetAddedAt() int64 {
	if x != nil {
		return x.AddedAt
	}
	return 0
}

var File_proto_store_proto protoreflect.FileDescriptor

var file_proto_store_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0xfb, 0x01, 0x0a, 0x0e, 0x50,
	0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x61, 0x75, 0x73, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d,
	0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x19, 0x0a,
	0x08, 0x61, 0x64, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x61, 0x64, 0x64, 0x65, 0x64, 0x41, 0x74, 0x22, 0x54, 0x0a, 0x0f, 0x50, 0x65, 0x72, 0x73,
	0x69, 0x73, 0x74, 0x65, 0x64, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x64, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x61, 0x64, 0x64, 0x65, 0x64, 0x41, 0x74, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_store_proto_rawDescOnce sync.Once
	file_proto_store_proto_rawDescData = file_proto_store_proto_rawDesc
)

func file_proto_store_proto_rawDescGZIP() []byte {
	file_proto_store_proto_rawDescOnce.Do(func() {
		file_proto_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_store_proto_rawDescData)
	})
	return file_proto_store_proto_rawDescData
}

var file_proto_store_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_store_proto_goTypes = []interface{}{
	(*PersistedEvent)(nil),  // 0: store.PersistedEvent
	(*PersistedStream)(nil), // 1: store.PersistedStream
}
var file_proto_store_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_store_proto_init() }
func file_proto_store_proto_init() {
	if File_proto_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistedEvent); i {
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
		file_proto_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistedStream); i {
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
			RawDescriptor: file_proto_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_store_proto_goTypes,
		DependencyIndexes: file_proto_store_proto_depIdxs,
		MessageInfos:      file_proto_store_proto_msgTypes,
	}.Build()
	File_proto_store_proto = out.File
	file_proto_store_proto_rawDesc = nil
	file_proto_store_proto_goTypes = nil
	file_proto_store_proto_depIdxs = nil
}
