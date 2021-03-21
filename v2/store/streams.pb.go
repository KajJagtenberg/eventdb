// Code generated by protoc-gen-go. DO NOT EDIT.
// source: streams.proto

package store

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	return fileDescriptor_c6bbf8af0ec331d6, []int{0}
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

type RecordedStream struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Events               [][]byte `protobuf:"bytes,2,rep,name=events,proto3" json:"events,omitempty"`
	AddedAt              int64    `protobuf:"varint,3,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecordedStream) Reset()         { *m = RecordedStream{} }
func (m *RecordedStream) String() string { return proto.CompactTextString(m) }
func (*RecordedStream) ProtoMessage()    {}
func (*RecordedStream) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{1}
}

func (m *RecordedStream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordedStream.Unmarshal(m, b)
}
func (m *RecordedStream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordedStream.Marshal(b, m, deterministic)
}
func (m *RecordedStream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordedStream.Merge(m, src)
}
func (m *RecordedStream) XXX_Size() int {
	return xxx_messageInfo_RecordedStream.Size(m)
}
func (m *RecordedStream) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordedStream.DiscardUnknown(m)
}

var xxx_messageInfo_RecordedStream proto.InternalMessageInfo

func (m *RecordedStream) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *RecordedStream) GetEvents() [][]byte {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *RecordedStream) GetAddedAt() int64 {
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
	return fileDescriptor_c6bbf8af0ec331d6, []int{2}
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
	return fileDescriptor_c6bbf8af0ec331d6, []int{2, 0}
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

type GetRequest struct {
	Stream               []byte   `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
	Version              uint32   `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Limit                uint32   `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{3}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetStream() []byte {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *GetRequest) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *GetRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type LogRequest struct {
	Offset               []byte   `protobuf:"bytes,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                uint32   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogRequest) Reset()         { *m = LogRequest{} }
func (m *LogRequest) String() string { return proto.CompactTextString(m) }
func (*LogRequest) ProtoMessage()    {}
func (*LogRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{4}
}

func (m *LogRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogRequest.Unmarshal(m, b)
}
func (m *LogRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogRequest.Marshal(b, m, deterministic)
}
func (m *LogRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogRequest.Merge(m, src)
}
func (m *LogRequest) XXX_Size() int {
	return xxx_messageInfo_LogRequest.Size(m)
}
func (m *LogRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogRequest proto.InternalMessageInfo

func (m *LogRequest) GetOffset() []byte {
	if m != nil {
		return m.Offset
	}
	return nil
}

func (m *LogRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type StreamCountRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamCountRequest) Reset()         { *m = StreamCountRequest{} }
func (m *StreamCountRequest) String() string { return proto.CompactTextString(m) }
func (*StreamCountRequest) ProtoMessage()    {}
func (*StreamCountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{5}
}

func (m *StreamCountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamCountRequest.Unmarshal(m, b)
}
func (m *StreamCountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamCountRequest.Marshal(b, m, deterministic)
}
func (m *StreamCountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamCountRequest.Merge(m, src)
}
func (m *StreamCountRequest) XXX_Size() int {
	return xxx_messageInfo_StreamCountRequest.Size(m)
}
func (m *StreamCountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamCountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamCountRequest proto.InternalMessageInfo

type StreamCountResponse struct {
	Count                uint64   `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamCountResponse) Reset()         { *m = StreamCountResponse{} }
func (m *StreamCountResponse) String() string { return proto.CompactTextString(m) }
func (*StreamCountResponse) ProtoMessage()    {}
func (*StreamCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{6}
}

func (m *StreamCountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamCountResponse.Unmarshal(m, b)
}
func (m *StreamCountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamCountResponse.Marshal(b, m, deterministic)
}
func (m *StreamCountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamCountResponse.Merge(m, src)
}
func (m *StreamCountResponse) XXX_Size() int {
	return xxx_messageInfo_StreamCountResponse.Size(m)
}
func (m *StreamCountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamCountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StreamCountResponse proto.InternalMessageInfo

func (m *StreamCountResponse) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type GetStreamsRequest struct {
	Skip                 uint32   `protobuf:"varint,1,opt,name=skip,proto3" json:"skip,omitempty"`
	Limit                uint32   `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetStreamsRequest) Reset()         { *m = GetStreamsRequest{} }
func (m *GetStreamsRequest) String() string { return proto.CompactTextString(m) }
func (*GetStreamsRequest) ProtoMessage()    {}
func (*GetStreamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{7}
}

func (m *GetStreamsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetStreamsRequest.Unmarshal(m, b)
}
func (m *GetStreamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetStreamsRequest.Marshal(b, m, deterministic)
}
func (m *GetStreamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetStreamsRequest.Merge(m, src)
}
func (m *GetStreamsRequest) XXX_Size() int {
	return xxx_messageInfo_GetStreamsRequest.Size(m)
}
func (m *GetStreamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetStreamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetStreamsRequest proto.InternalMessageInfo

func (m *GetStreamsRequest) GetSkip() uint32 {
	if m != nil {
		return m.Skip
	}
	return 0
}

func (m *GetStreamsRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type Stream struct {
	Id                   []byte   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AddedAt              int64    `protobuf:"varint,2,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stream) Reset()         { *m = Stream{} }
func (m *Stream) String() string { return proto.CompactTextString(m) }
func (*Stream) ProtoMessage()    {}
func (*Stream) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{8}
}

func (m *Stream) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stream.Unmarshal(m, b)
}
func (m *Stream) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stream.Marshal(b, m, deterministic)
}
func (m *Stream) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stream.Merge(m, src)
}
func (m *Stream) XXX_Size() int {
	return xxx_messageInfo_Stream.Size(m)
}
func (m *Stream) XXX_DiscardUnknown() {
	xxx_messageInfo_Stream.DiscardUnknown(m)
}

var xxx_messageInfo_Stream proto.InternalMessageInfo

func (m *Stream) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Stream) GetAddedAt() int64 {
	if m != nil {
		return m.AddedAt
	}
	return 0
}

type EventCountRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventCountRequest) Reset()         { *m = EventCountRequest{} }
func (m *EventCountRequest) String() string { return proto.CompactTextString(m) }
func (*EventCountRequest) ProtoMessage()    {}
func (*EventCountRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{9}
}

func (m *EventCountRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventCountRequest.Unmarshal(m, b)
}
func (m *EventCountRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventCountRequest.Marshal(b, m, deterministic)
}
func (m *EventCountRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventCountRequest.Merge(m, src)
}
func (m *EventCountRequest) XXX_Size() int {
	return xxx_messageInfo_EventCountRequest.Size(m)
}
func (m *EventCountRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EventCountRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EventCountRequest proto.InternalMessageInfo

type EventCountResponse struct {
	Count                uint64   `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventCountResponse) Reset()         { *m = EventCountResponse{} }
func (m *EventCountResponse) String() string { return proto.CompactTextString(m) }
func (*EventCountResponse) ProtoMessage()    {}
func (*EventCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6bbf8af0ec331d6, []int{10}
}

func (m *EventCountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventCountResponse.Unmarshal(m, b)
}
func (m *EventCountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventCountResponse.Marshal(b, m, deterministic)
}
func (m *EventCountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventCountResponse.Merge(m, src)
}
func (m *EventCountResponse) XXX_Size() int {
	return xxx_messageInfo_EventCountResponse.Size(m)
}
func (m *EventCountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EventCountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EventCountResponse proto.InternalMessageInfo

func (m *EventCountResponse) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func init() {
	proto.RegisterType((*RecordedEvent)(nil), "store.RecordedEvent")
	proto.RegisterType((*RecordedStream)(nil), "store.RecordedStream")
	proto.RegisterType((*AddRequest)(nil), "store.AddRequest")
	proto.RegisterType((*AddRequest_Event)(nil), "store.AddRequest.Event")
	proto.RegisterType((*GetRequest)(nil), "store.GetRequest")
	proto.RegisterType((*LogRequest)(nil), "store.LogRequest")
	proto.RegisterType((*StreamCountRequest)(nil), "store.StreamCountRequest")
	proto.RegisterType((*StreamCountResponse)(nil), "store.StreamCountResponse")
	proto.RegisterType((*GetStreamsRequest)(nil), "store.GetStreamsRequest")
	proto.RegisterType((*Stream)(nil), "store.Stream")
	proto.RegisterType((*EventCountRequest)(nil), "store.EventCountRequest")
	proto.RegisterType((*EventCountResponse)(nil), "store.EventCountResponse")
}

func init() { proto.RegisterFile("streams.proto", fileDescriptor_c6bbf8af0ec331d6) }

var fileDescriptor_c6bbf8af0ec331d6 = []byte{
	// 536 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0x5e, 0x3b, 0x49, 0xd3, 0x9d, 0x6e, 0x2a, 0xd5, 0x5b, 0x81, 0x9b, 0x53, 0x88, 0x84, 0x14,
	0x81, 0x54, 0x50, 0xf7, 0x04, 0x12, 0x87, 0xb2, 0x82, 0x6a, 0xa5, 0x3d, 0x65, 0x39, 0x71, 0x59,
	0x85, 0xda, 0xbb, 0x8a, 0xd8, 0xc6, 0x25, 0x76, 0x57, 0xe2, 0x09, 0x78, 0x02, 0x1e, 0x93, 0x17,
	0xe0, 0x84, 0x62, 0xe7, 0x57, 0x4d, 0x0b, 0xe2, 0xe6, 0xf9, 0x3c, 0xf3, 0x79, 0xe6, 0xfb, 0x46,
	0x06, 0x4f, 0xaa, 0x9c, 0x27, 0x1b, 0x39, 0xdf, 0xe6, 0x42, 0x09, 0xe2, 0x48, 0x25, 0x72, 0x1e,
	0xfe, 0x46, 0xe0, 0xc5, 0x7c, 0x2d, 0x72, 0xc6, 0xd9, 0x87, 0x47, 0x9e, 0x29, 0x32, 0x06, 0x9c,
	0x32, 0x8a, 0x02, 0x14, 0x9d, 0xc5, 0x38, 0x65, 0xe4, 0x09, 0x0c, 0x4c, 0x25, 0xc5, 0x1a, 0x2b,
	0x23, 0x42, 0xc1, 0x7d, 0xe4, 0xb9, 0x4c, 0x45, 0x46, 0xad, 0x00, 0x45, 0x5e, 0x5c, 0x85, 0x84,
	0x80, 0xad, 0xbe, 0x6f, 0x39, 0xb5, 0x03, 0x14, 0x9d, 0xc6, 0xfa, 0x5c, 0x60, 0x2c, 0x51, 0x09,
	0x75, 0x34, 0x87, 0x3e, 0x13, 0x1f, 0x86, 0x1b, 0xae, 0x12, 0x8d, 0x0f, 0x34, 0x5e, 0xc7, 0xe4,
	0x19, 0x9c, 0xad, 0x93, 0x9d, 0x4c, 0x54, 0x2a, 0xb2, 0xdb, 0x94, 0x51, 0x57, 0xdf, 0x8f, 0x6a,
	0xec, 0x8a, 0x91, 0xe7, 0x30, 0x5e, 0x8b, 0x3c, 0xe7, 0x0f, 0x75, 0xd2, 0x50, 0x27, 0x79, 0x2d,
	0xf4, 0x8a, 0x91, 0x19, 0x0c, 0x13, 0xc6, 0x38, 0xbb, 0x4d, 0x14, 0x3d, 0x0d, 0x50, 0x64, 0xc5,
	0xae, 0x8e, 0x97, 0x2a, 0xbc, 0x81, 0x71, 0x35, 0xfb, 0x8d, 0x19, 0xaa, 0x67, 0x78, 0x5e, 0xa8,
	0x22, 0x29, 0x0e, 0xac, 0x62, 0x78, 0x13, 0x75, 0x48, 0xad, 0x2e, 0xe9, 0x0f, 0x0c, 0xb0, 0x64,
	0x2c, 0xe6, 0xdf, 0x76, 0x5c, 0xaa, 0x96, 0x7c, 0xe8, 0x90, 0x7c, 0xb8, 0x2b, 0xdf, 0xab, 0xfa,
	0x4d, 0x2b, 0xb0, 0xa2, 0xd1, 0xe2, 0xe9, 0x5c, 0x5b, 0x35, 0x6f, 0x48, 0xe7, 0xda, 0xa9, 0xaa,
	0x19, 0xff, 0x27, 0x02, 0xc7, 0x78, 0x57, 0x29, 0x8f, 0x7a, 0x94, 0xc7, 0x07, 0x94, 0xb7, 0xfe,
	0xa2, 0xbc, 0xfd, 0x2f, 0xca, 0x3b, 0x3d, 0xca, 0x87, 0x9f, 0x00, 0x56, 0x5c, 0xfd, 0xbf, 0x10,
	0x53, 0x70, 0x1e, 0xd2, 0x4d, 0xaa, 0xca, 0xfd, 0x32, 0x41, 0xf8, 0x16, 0xe0, 0x5a, 0xdc, 0xb7,
	0x58, 0xc5, 0xdd, 0x9d, 0xe4, 0xaa, 0x62, 0x35, 0x51, 0x53, 0x8b, 0xdb, 0xb5, 0x53, 0x20, 0xc6,
	0xe8, 0x4b, 0xb1, 0xcb, 0xaa, 0xce, 0xc2, 0x97, 0x70, 0xde, 0x41, 0xe5, 0x56, 0x64, 0x92, 0x17,
	0x14, 0xeb, 0x02, 0xd0, 0xcc, 0x76, 0x6c, 0x82, 0xf0, 0x1d, 0x4c, 0x56, 0x5c, 0x99, 0x7c, 0x59,
	0x75, 0x41, 0xc0, 0x96, 0x5f, 0xd3, 0xad, 0xce, 0xf4, 0x62, 0x7d, 0x3e, 0xd0, 0xc1, 0x05, 0x0c,
	0x0e, 0xac, 0x5a, 0x7b, 0xa5, 0x70, 0x77, 0xa5, 0xce, 0x61, 0xa2, 0xfd, 0xed, 0x74, 0xfd, 0x02,
	0x48, 0x1b, 0x3c, 0xd6, 0xf4, 0xe2, 0x17, 0x06, 0xb7, 0x6c, 0x99, 0x2c, 0xc0, 0x5a, 0x32, 0x46,
	0x26, 0x7b, 0x5b, 0xe5, 0x4f, 0x4b, 0xa8, 0xf3, 0x1f, 0x84, 0x27, 0xaf, 0x51, 0x51, 0xb3, 0xe2,
	0xaa, 0xae, 0x69, 0x5c, 0x3d, 0x5e, 0x73, 0x2d, 0xee, 0xeb, 0x9a, 0xc6, 0xb3, 0x23, 0x35, 0x1f,
	0x61, 0xd4, 0x72, 0x82, 0xcc, 0xca, 0xc4, 0x7d, 0xcf, 0x7c, 0xbf, 0xef, 0xca, 0x68, 0x10, 0x9e,
	0x90, 0x4b, 0x80, 0x46, 0x1b, 0x42, 0xcb, 0xdc, 0x3d, 0x0d, 0xfd, 0x59, 0xcf, 0x4d, 0x4d, 0xf2,
	0x46, 0xaf, 0x6f, 0x25, 0x1b, 0x6d, 0x66, 0xef, 0x9a, 0xef, 0x7b, 0x9d, 0x56, 0x8a, 0x39, 0xde,
	0xbb, 0x9f, 0xcd, 0xf7, 0xfa, 0x65, 0xa0, 0x3f, 0xdb, 0x8b, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xdb, 0x86, 0x1f, 0x24, 0x7d, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StreamsClient is the client API for Streams service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamsClient interface {
	Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (Streams_AddClient, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Streams_GetClient, error)
	Log(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (Streams_LogClient, error)
	StreamCount(ctx context.Context, in *StreamCountRequest, opts ...grpc.CallOption) (*StreamCountResponse, error)
	EventCount(ctx context.Context, in *EventCountRequest, opts ...grpc.CallOption) (*EventCountResponse, error)
	GetStreams(ctx context.Context, in *GetStreamsRequest, opts ...grpc.CallOption) (Streams_GetStreamsClient, error)
}

type streamsClient struct {
	cc *grpc.ClientConn
}

func NewStreamsClient(cc *grpc.ClientConn) StreamsClient {
	return &streamsClient{cc}
}

func (c *streamsClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (Streams_AddClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Streams_serviceDesc.Streams[0], "/store.Streams/Add", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamsAddClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Streams_AddClient interface {
	Recv() (*RecordedEvent, error)
	grpc.ClientStream
}

type streamsAddClient struct {
	grpc.ClientStream
}

func (x *streamsAddClient) Recv() (*RecordedEvent, error) {
	m := new(RecordedEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamsClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Streams_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Streams_serviceDesc.Streams[1], "/store.Streams/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamsGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Streams_GetClient interface {
	Recv() (*RecordedEvent, error)
	grpc.ClientStream
}

type streamsGetClient struct {
	grpc.ClientStream
}

func (x *streamsGetClient) Recv() (*RecordedEvent, error) {
	m := new(RecordedEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamsClient) Log(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (Streams_LogClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Streams_serviceDesc.Streams[2], "/store.Streams/Log", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamsLogClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Streams_LogClient interface {
	Recv() (*RecordedEvent, error)
	grpc.ClientStream
}

type streamsLogClient struct {
	grpc.ClientStream
}

func (x *streamsLogClient) Recv() (*RecordedEvent, error) {
	m := new(RecordedEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamsClient) StreamCount(ctx context.Context, in *StreamCountRequest, opts ...grpc.CallOption) (*StreamCountResponse, error) {
	out := new(StreamCountResponse)
	err := c.cc.Invoke(ctx, "/store.Streams/StreamCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamsClient) EventCount(ctx context.Context, in *EventCountRequest, opts ...grpc.CallOption) (*EventCountResponse, error) {
	out := new(EventCountResponse)
	err := c.cc.Invoke(ctx, "/store.Streams/EventCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamsClient) GetStreams(ctx context.Context, in *GetStreamsRequest, opts ...grpc.CallOption) (Streams_GetStreamsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Streams_serviceDesc.Streams[3], "/store.Streams/GetStreams", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamsGetStreamsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Streams_GetStreamsClient interface {
	Recv() (*Stream, error)
	grpc.ClientStream
}

type streamsGetStreamsClient struct {
	grpc.ClientStream
}

func (x *streamsGetStreamsClient) Recv() (*Stream, error) {
	m := new(Stream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamsServer is the server API for Streams service.
type StreamsServer interface {
	Add(*AddRequest, Streams_AddServer) error
	Get(*GetRequest, Streams_GetServer) error
	Log(*LogRequest, Streams_LogServer) error
	StreamCount(context.Context, *StreamCountRequest) (*StreamCountResponse, error)
	EventCount(context.Context, *EventCountRequest) (*EventCountResponse, error)
	GetStreams(*GetStreamsRequest, Streams_GetStreamsServer) error
}

// UnimplementedStreamsServer can be embedded to have forward compatible implementations.
type UnimplementedStreamsServer struct {
}

func (*UnimplementedStreamsServer) Add(req *AddRequest, srv Streams_AddServer) error {
	return status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (*UnimplementedStreamsServer) Get(req *GetRequest, srv Streams_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedStreamsServer) Log(req *LogRequest, srv Streams_LogServer) error {
	return status.Errorf(codes.Unimplemented, "method Log not implemented")
}
func (*UnimplementedStreamsServer) StreamCount(ctx context.Context, req *StreamCountRequest) (*StreamCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StreamCount not implemented")
}
func (*UnimplementedStreamsServer) EventCount(ctx context.Context, req *EventCountRequest) (*EventCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EventCount not implemented")
}
func (*UnimplementedStreamsServer) GetStreams(req *GetStreamsRequest, srv Streams_GetStreamsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStreams not implemented")
}

func RegisterStreamsServer(s *grpc.Server, srv StreamsServer) {
	s.RegisterService(&_Streams_serviceDesc, srv)
}

func _Streams_Add_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AddRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamsServer).Add(m, &streamsAddServer{stream})
}

type Streams_AddServer interface {
	Send(*RecordedEvent) error
	grpc.ServerStream
}

type streamsAddServer struct {
	grpc.ServerStream
}

func (x *streamsAddServer) Send(m *RecordedEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Streams_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamsServer).Get(m, &streamsGetServer{stream})
}

type Streams_GetServer interface {
	Send(*RecordedEvent) error
	grpc.ServerStream
}

type streamsGetServer struct {
	grpc.ServerStream
}

func (x *streamsGetServer) Send(m *RecordedEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Streams_Log_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LogRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamsServer).Log(m, &streamsLogServer{stream})
}

type Streams_LogServer interface {
	Send(*RecordedEvent) error
	grpc.ServerStream
}

type streamsLogServer struct {
	grpc.ServerStream
}

func (x *streamsLogServer) Send(m *RecordedEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Streams_StreamCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamsServer).StreamCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.Streams/StreamCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamsServer).StreamCount(ctx, req.(*StreamCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Streams_EventCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamsServer).EventCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/store.Streams/EventCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamsServer).EventCount(ctx, req.(*EventCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Streams_GetStreams_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetStreamsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamsServer).GetStreams(m, &streamsGetStreamsServer{stream})
}

type Streams_GetStreamsServer interface {
	Send(*Stream) error
	grpc.ServerStream
}

type streamsGetStreamsServer struct {
	grpc.ServerStream
}

func (x *streamsGetStreamsServer) Send(m *Stream) error {
	return x.ServerStream.SendMsg(m)
}

var _Streams_serviceDesc = grpc.ServiceDesc{
	ServiceName: "store.Streams",
	HandlerType: (*StreamsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StreamCount",
			Handler:    _Streams_StreamCount_Handler,
		},
		{
			MethodName: "EventCount",
			Handler:    _Streams_EventCount_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Add",
			Handler:       _Streams_Add_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Get",
			Handler:       _Streams_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Log",
			Handler:       _Streams_Log_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStreams",
			Handler:       _Streams_GetStreams_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "streams.proto",
}