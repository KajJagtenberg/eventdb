// Code generated by protoc-gen-go. DO NOT EDIT.
// source: store.proto

package store

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

type PersistedEvent struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Stream               []byte   `protobuf:"bytes,2,opt,name=stream,proto3" json:"stream,omitempty"`
	Version              uint32   `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Type                 string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Data                 []byte   `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Metadata             []byte   `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CausationId          []byte   `protobuf:"bytes,7,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId        []byte   `protobuf:"bytes,8,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	AddedAt              int64    `protobuf:"varint,9,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PersistedEvent) Reset()         { *m = PersistedEvent{} }
func (m *PersistedEvent) String() string { return proto.CompactTextString(m) }
func (*PersistedEvent) ProtoMessage()    {}
func (*PersistedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_98bbca36ef968dfc, []int{0}
}

func (m *PersistedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PersistedEvent.Unmarshal(m, b)
}
func (m *PersistedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PersistedEvent.Marshal(b, m, deterministic)
}
func (m *PersistedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PersistedEvent.Merge(m, src)
}
func (m *PersistedEvent) XXX_Size() int {
	return xxx_messageInfo_PersistedEvent.Size(m)
}
func (m *PersistedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_PersistedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_PersistedEvent proto.InternalMessageInfo

func (m *PersistedEvent) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *PersistedEvent) GetStream() []byte {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *PersistedEvent) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PersistedEvent) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *PersistedEvent) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PersistedEvent) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *PersistedEvent) GetCausationId() []byte {
	if m != nil {
		return m.CausationId
	}
	return nil
}

func (m *PersistedEvent) GetCorrelationId() []byte {
	if m != nil {
		return m.CorrelationId
	}
	return nil
}

func (m *PersistedEvent) GetAddedAt() int64 {
	if m != nil {
		return m.AddedAt
	}
	return 0
}

type PersistedStream struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Events               [][]byte `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
	AddedAt              int64    `protobuf:"varint,3,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PersistedStream) Reset()         { *m = PersistedStream{} }
func (m *PersistedStream) String() string { return proto.CompactTextString(m) }
func (*PersistedStream) ProtoMessage()    {}
func (*PersistedStream) Descriptor() ([]byte, []int) {
	return fileDescriptor_98bbca36ef968dfc, []int{1}
}

func (m *PersistedStream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PersistedStream.Unmarshal(m, b)
}
func (m *PersistedStream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PersistedStream.Marshal(b, m, deterministic)
}
func (m *PersistedStream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PersistedStream.Merge(m, src)
}
func (m *PersistedStream) XXX_Size() int {
	return xxx_messageInfo_PersistedStream.Size(m)
}
func (m *PersistedStream) XXX_DiscardUnknown() {
	xxx_messageInfo_PersistedStream.DiscardUnknown(m)
}

var xxx_messageInfo_PersistedStream proto.InternalMessageInfo

func (m *PersistedStream) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *PersistedStream) GetEvents() [][]byte {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *PersistedStream) GetAddedAt() int64 {
	if m != nil {
		return m.AddedAt
	}
	return 0
}

func init() {
	proto.RegisterType((*PersistedEvent)(nil), "store.PersistedEvent")
	proto.RegisterType((*PersistedStream)(nil), "store.PersistedStream")
}

func init() { proto.RegisterFile("store.proto", fileDescriptor_98bbca36ef968dfc) }

var fileDescriptor_98bbca36ef968dfc = []byte{
	// 250 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0x86, 0xd9, 0xa4, 0x4d, 0xd2, 0x69, 0x9b, 0x0f, 0xe6, 0xf0, 0xb1, 0x7a, 0x8a, 0x05, 0x21,
	0x27, 0x2f, 0xfe, 0x02, 0x05, 0x0f, 0xbd, 0x49, 0xf4, 0xe4, 0xa5, 0xac, 0x9d, 0x39, 0x2c, 0xd8,
	0x6c, 0xd9, 0x1d, 0x0b, 0xfe, 0x76, 0x2f, 0x92, 0x89, 0x06, 0x0b, 0xde, 0xe6, 0x79, 0x76, 0x76,
	0xe0, 0x7d, 0x61, 0x99, 0x24, 0x44, 0xbe, 0x39, 0xc6, 0x20, 0x01, 0xe7, 0x0a, 0x9b, 0x4f, 0x03,
	0xf5, 0x23, 0xc7, 0xe4, 0x93, 0x30, 0x3d, 0x9c, 0xb8, 0x17, 0xac, 0x21, 0xf3, 0x64, 0x4d, 0x63,
	0xda, 0x55, 0x97, 0x79, 0xc2, 0xff, 0x50, 0x24, 0x89, 0xec, 0x0e, 0x36, 0x53, 0xf7, 0x4d, 0x68,
	0xa1, 0x3c, 0x0d, 0x3f, 0x43, 0x6f, 0xf3, 0xc6, 0xb4, 0xeb, 0xee, 0x07, 0x11, 0x61, 0x26, 0x1f,
	0x47, 0xb6, 0xb3, 0xc6, 0xb4, 0x8b, 0x4e, 0xe7, 0xc1, 0x91, 0x13, 0x67, 0xe7, 0x7a, 0x43, 0x67,
	0xbc, 0x84, 0xea, 0xc0, 0xe2, 0xd4, 0x17, 0xea, 0x27, 0xc6, 0x2b, 0x58, 0xed, 0xdd, 0x7b, 0x72,
	0xe2, 0x43, 0xbf, 0xf3, 0x64, 0x4b, 0x7d, 0x5f, 0x4e, 0x6e, 0x4b, 0x78, 0x0d, 0xf5, 0x3e, 0xc4,
	0xc8, 0x6f, 0xd3, 0x52, 0xa5, 0x4b, 0xeb, 0x5f, 0x76, 0x4b, 0x78, 0x01, 0x95, 0x23, 0x62, 0xda,
	0x39, 0xb1, 0x8b, 0xc6, 0xb4, 0x79, 0x57, 0x2a, 0xdf, 0xc9, 0xe6, 0x19, 0xfe, 0x4d, 0xe1, 0x9f,
	0xc6, 0x54, 0x7f, 0xa4, 0xe7, 0xa1, 0x96, 0x64, 0xb3, 0x26, 0x1f, 0xd2, 0x8f, 0x74, 0x76, 0x35,
	0x3f, 0xbb, 0x7a, 0x5f, 0xbe, 0x8c, 0xe5, 0xbe, 0x16, 0x5a, 0xf5, 0xed, 0x57, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x18, 0x24, 0x0c, 0x77, 0x79, 0x01, 0x00, 0x00,
}
