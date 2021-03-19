// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/store.proto

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

type RecordedEvent struct {
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

func (m *RecordedEvent) Reset()         { *m = RecordedEvent{} }
func (m *RecordedEvent) String() string { return proto.CompactTextString(m) }
func (*RecordedEvent) ProtoMessage()    {}
func (*RecordedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f39e0831e20e871, []int{0}
}

func (m *RecordedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordedEvent.Unmarshal(m, b)
}
func (m *RecordedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordedEvent.Marshal(b, m, deterministic)
}
func (m *RecordedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordedEvent.Merge(m, src)
}
func (m *RecordedEvent) XXX_Size() int {
	return xxx_messageInfo_RecordedEvent.Size(m)
}
func (m *RecordedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_RecordedEvent proto.InternalMessageInfo

func (m *RecordedEvent) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *RecordedEvent) GetStream() []byte {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *RecordedEvent) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *RecordedEvent) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *RecordedEvent) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *RecordedEvent) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *RecordedEvent) GetCausationId() []byte {
	if m != nil {
		return m.CausationId
	}
	return nil
}

func (m *RecordedEvent) GetCorrelationId() []byte {
	if m != nil {
		return m.CorrelationId
	}
	return nil
}

func (m *RecordedEvent) GetAddedAt() int64 {
	if m != nil {
		return m.AddedAt
	}
	return 0
}

type AddRequest struct {
	Stream               []byte              `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
	Version              uint32              `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Events               []*AddRequest_Event `protobuf:"bytes,3,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}
func (*AddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f39e0831e20e871, []int{1}
}

func (m *AddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRequest.Unmarshal(m, b)
}
func (m *AddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRequest.Marshal(b, m, deterministic)
}
func (m *AddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRequest.Merge(m, src)
}
func (m *AddRequest) XXX_Size() int {
	return xxx_messageInfo_AddRequest.Size(m)
}
func (m *AddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRequest proto.InternalMessageInfo

func (m *AddRequest) GetStream() []byte {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *AddRequest) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *AddRequest) GetEvents() []*AddRequest_Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type AddRequest_Event struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Metadata             []byte   `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CausationId          []byte   `protobuf:"bytes,4,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId        []byte   `protobuf:"bytes,5,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRequest_Event) Reset()         { *m = AddRequest_Event{} }
func (m *AddRequest_Event) String() string { return proto.CompactTextString(m) }
func (*AddRequest_Event) ProtoMessage()    {}
func (*AddRequest_Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f39e0831e20e871, []int{1, 0}
}

func (m *AddRequest_Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRequest_Event.Unmarshal(m, b)
}
func (m *AddRequest_Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRequest_Event.Marshal(b, m, deterministic)
}
func (m *AddRequest_Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRequest_Event.Merge(m, src)
}
func (m *AddRequest_Event) XXX_Size() int {
	return xxx_messageInfo_AddRequest_Event.Size(m)
}
func (m *AddRequest_Event) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRequest_Event.DiscardUnknown(m)
}

var xxx_messageInfo_AddRequest_Event proto.InternalMessageInfo

func (m *AddRequest_Event) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *AddRequest_Event) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *AddRequest_Event) GetMetadata() []byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *AddRequest_Event) GetCausationId() []byte {
	if m != nil {
		return m.CausationId
	}
	return nil
}

func (m *AddRequest_Event) GetCorrelationId() []byte {
	if m != nil {
		return m.CorrelationId
	}
	return nil
}

type AddResponse struct {
	Events               []*RecordedEvent `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *AddResponse) Reset()         { *m = AddResponse{} }
func (m *AddResponse) String() string { return proto.CompactTextString(m) }
func (*AddResponse) ProtoMessage()    {}
func (*AddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f39e0831e20e871, []int{2}
}

func (m *AddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddResponse.Unmarshal(m, b)
}
func (m *AddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddResponse.Marshal(b, m, deterministic)
}
func (m *AddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddResponse.Merge(m, src)
}
func (m *AddResponse) XXX_Size() int {
	return xxx_messageInfo_AddResponse.Size(m)
}
func (m *AddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddResponse proto.InternalMessageInfo

func (m *AddResponse) GetEvents() []*RecordedEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterType((*RecordedEvent)(nil), "store.RecordedEvent")
	proto.RegisterType((*AddRequest)(nil), "store.AddRequest")
	proto.RegisterType((*AddRequest_Event)(nil), "store.AddRequest.Event")
	proto.RegisterType((*AddResponse)(nil), "store.AddResponse")
}

func init() { proto.RegisterFile("proto/store.proto", fileDescriptor_3f39e0831e20e871) }

var fileDescriptor_3f39e0831e20e871 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcd, 0x4e, 0xf2, 0x40,
	0x14, 0x86, 0xbf, 0x69, 0x69, 0x81, 0xc3, 0x4f, 0xc2, 0xe4, 0x8b, 0x8e, 0xac, 0x6a, 0x13, 0x93,
	0x2e, 0x4c, 0x49, 0x70, 0xa9, 0x1b, 0x16, 0x2e, 0xd8, 0xd6, 0x0b, 0x20, 0x23, 0xe7, 0x2c, 0x9a,
	0x48, 0x07, 0x67, 0x06, 0x12, 0xaf, 0xc0, 0x2b, 0xf0, 0x66, 0x5d, 0x19, 0x0e, 0x50, 0x21, 0x42,
	0x74, 0x77, 0xde, 0xa7, 0x27, 0xd3, 0x79, 0x9f, 0x0c, 0x0c, 0x96, 0xd6, 0x78, 0x33, 0x72, 0xde,
	0x58, 0xca, 0x79, 0x96, 0x11, 0x87, 0xf4, 0x53, 0x40, 0xaf, 0xa0, 0xb9, 0xb1, 0x48, 0xf8, 0xb8,
	0xa6, 0xca, 0xcb, 0x3e, 0x04, 0x25, 0x2a, 0x91, 0x88, 0xac, 0x5b, 0x04, 0x25, 0xca, 0x0b, 0x88,
	0x9d, 0xb7, 0xa4, 0x17, 0x2a, 0x60, 0xb6, 0x4b, 0x52, 0x41, 0x73, 0x4d, 0xd6, 0x95, 0xa6, 0x52,
	0x61, 0x22, 0xb2, 0x5e, 0xb1, 0x8f, 0x52, 0x42, 0xc3, 0xbf, 0x2d, 0x49, 0x35, 0x12, 0x91, 0xb5,
	0x0b, 0x9e, 0x37, 0x0c, 0xb5, 0xd7, 0x2a, 0xe2, 0x33, 0x78, 0x96, 0x43, 0x68, 0x2d, 0xc8, 0x6b,
	0xe6, 0x31, 0xf3, 0x3a, 0xcb, 0x6b, 0xe8, 0xce, 0xf5, 0xca, 0x69, 0x5f, 0x9a, 0x6a, 0x56, 0xa2,
	0x6a, 0xf2, 0xf7, 0x4e, 0xcd, 0xa6, 0x28, 0x6f, 0xa0, 0x3f, 0x37, 0xd6, 0xd2, 0x4b, 0xbd, 0xd4,
	0xe2, 0xa5, 0xde, 0x01, 0x9d, 0xa2, 0xbc, 0x82, 0x96, 0x46, 0x24, 0x9c, 0x69, 0xaf, 0xda, 0x89,
	0xc8, 0xc2, 0xa2, 0xc9, 0x79, 0xe2, 0xd3, 0xf7, 0x00, 0x60, 0x82, 0x58, 0xd0, 0xeb, 0x8a, 0x9c,
	0x3f, 0x68, 0x2a, 0xce, 0x35, 0x0d, 0x8e, 0x9b, 0x8e, 0x20, 0xa6, 0x8d, 0x34, 0xa7, 0xc2, 0x24,
	0xcc, 0x3a, 0xe3, 0xcb, 0x7c, 0xab, 0xf8, 0xfb, 0xd0, 0x9c, 0xa5, 0x16, 0xbb, 0xb5, 0xe1, 0x87,
	0x80, 0x68, 0xab, 0x79, 0x2f, 0x49, 0x9c, 0x90, 0x14, 0x9c, 0x91, 0x14, 0xfe, 0x22, 0xa9, 0xf1,
	0x17, 0x49, 0xd1, 0x09, 0x49, 0xe9, 0x3d, 0x74, 0xf8, 0xce, 0x6e, 0x69, 0x2a, 0x47, 0xf2, 0xb6,
	0xee, 0x25, 0xb8, 0xd7, 0xff, 0x5d, 0xaf, 0xa3, 0x97, 0xb2, 0x2f, 0x35, 0x7e, 0x00, 0x60, 0xf0,
	0xb4, 0xd9, 0x91, 0x39, 0x84, 0x13, 0x44, 0x39, 0xf8, 0xa1, 0x62, 0x28, 0x0f, 0xd1, 0xf6, 0x4f,
	0xe9, 0xbf, 0xe7, 0x98, 0xdf, 0xe3, 0xdd, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x40, 0xff, 0x01,
	0x8e, 0xa4, 0x02, 0x00, 0x00,
}
