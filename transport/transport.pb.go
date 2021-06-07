// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: proto/transport.proto

package transport

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

type EventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*EventResponse_Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *EventResponse) Reset() {
	*x = EventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventResponse) ProtoMessage() {}

func (x *EventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventResponse.ProtoReflect.Descriptor instead.
func (*EventResponse) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{0}
}

func (x *EventResponse) GetEvents() []*EventResponse_Event {
	if x != nil {
		return x.Events
	}
	return nil
}

type AddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stream  string                  `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
	Version uint32                  `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Events  []*AddRequest_EventData `protobuf:"bytes,3,rep,name=events,proto3" json:"events,omitempty"`
}

func (x *AddRequest) Reset() {
	*x = AddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest) ProtoMessage() {}

func (x *AddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRequest.ProtoReflect.Descriptor instead.
func (*AddRequest) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{1}
}

func (x *AddRequest) GetStream() string {
	if x != nil {
		return x.Stream
	}
	return ""
}

func (x *AddRequest) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *AddRequest) GetEvents() []*AddRequest_EventData {
	if x != nil {
		return x.Events
	}
	return nil
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stream  string `protobuf:"bytes,1,opt,name=stream,proto3" json:"stream,omitempty"`
	Version uint32 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	Limit   uint32 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetStream() string {
	if x != nil {
		return x.Stream
	}
	return ""
}

func (x *GetRequest) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *GetRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type GetAllRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset string `protobuf:"bytes,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  uint32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *GetAllRequest) Reset() {
	*x = GetAllRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAllRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllRequest) ProtoMessage() {}

func (x *GetAllRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllRequest.ProtoReflect.Descriptor instead.
func (*GetAllRequest) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{3}
}

func (x *GetAllRequest) GetOffset() string {
	if x != nil {
		return x.Offset
	}
	return ""
}

func (x *GetAllRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type EventResponse_Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Stream        string `protobuf:"bytes,2,opt,name=stream,proto3" json:"stream,omitempty"`
	Version       uint32 `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Type          string `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Data          []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	Metadata      []byte `protobuf:"bytes,6,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CausationId   string `protobuf:"bytes,7,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId string `protobuf:"bytes,8,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	AddedAt       int64  `protobuf:"varint,9,opt,name=added_at,json=addedAt,proto3" json:"added_at,omitempty"`
}

func (x *EventResponse_Event) Reset() {
	*x = EventResponse_Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventResponse_Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventResponse_Event) ProtoMessage() {}

func (x *EventResponse_Event) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventResponse_Event.ProtoReflect.Descriptor instead.
func (*EventResponse_Event) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{0, 0}
}

func (x *EventResponse_Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EventResponse_Event) GetStream() string {
	if x != nil {
		return x.Stream
	}
	return ""
}

func (x *EventResponse_Event) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *EventResponse_Event) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *EventResponse_Event) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *EventResponse_Event) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *EventResponse_Event) GetCausationId() string {
	if x != nil {
		return x.CausationId
	}
	return ""
}

func (x *EventResponse_Event) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *EventResponse_Event) GetAddedAt() int64 {
	if x != nil {
		return x.AddedAt
	}
	return 0
}

type AddRequest_EventData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Data          []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Metadata      []byte `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	CausationId   string `protobuf:"bytes,4,opt,name=causation_id,json=causationId,proto3" json:"causation_id,omitempty"`
	CorrelationId string `protobuf:"bytes,5,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
}

func (x *AddRequest_EventData) Reset() {
	*x = AddRequest_EventData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_transport_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddRequest_EventData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddRequest_EventData) ProtoMessage() {}

func (x *AddRequest_EventData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_transport_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddRequest_EventData.ProtoReflect.Descriptor instead.
func (*AddRequest_EventData) Descriptor() ([]byte, []int) {
	return file_proto_transport_proto_rawDescGZIP(), []int{1, 0}
}

func (x *AddRequest_EventData) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *AddRequest_EventData) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AddRequest_EventData) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *AddRequest_EventData) GetCausationId() string {
	if x != nil {
		return x.CausationId
	}
	return ""
}

func (x *AddRequest_EventData) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

var File_proto_transport_proto protoreflect.FileDescriptor

var file_proto_transport_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f,
	0x72, 0x74, 0x22, 0xbc, 0x02, 0x0a, 0x0d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0xf2, 0x01, 0x0a,
	0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x0c,
	0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12,
	0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x64, 0x64, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x61, 0x64, 0x64, 0x65, 0x64, 0x41,
	0x74, 0x22, 0x93, 0x02, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x99, 0x01, 0x0a, 0x09,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a,
	0x0c, 0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x61, 0x75, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x54, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x3d, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x32, 0xc7, 0x01, 0x0a,
	0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x38, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x15, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x03,
	0x47, 0x65, 0x74, 0x12, 0x15, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c,
	0x12, 0x18, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x2e, 0x2f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_transport_proto_rawDescOnce sync.Once
	file_proto_transport_proto_rawDescData = file_proto_transport_proto_rawDesc
)

func file_proto_transport_proto_rawDescGZIP() []byte {
	file_proto_transport_proto_rawDescOnce.Do(func() {
		file_proto_transport_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_transport_proto_rawDescData)
	})
	return file_proto_transport_proto_rawDescData
}

var file_proto_transport_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_transport_proto_goTypes = []interface{}{
	(*EventResponse)(nil),        // 0: transport.EventResponse
	(*AddRequest)(nil),           // 1: transport.AddRequest
	(*GetRequest)(nil),           // 2: transport.GetRequest
	(*GetAllRequest)(nil),        // 3: transport.GetAllRequest
	(*EventResponse_Event)(nil),  // 4: transport.EventResponse.Event
	(*AddRequest_EventData)(nil), // 5: transport.AddRequest.EventData
}
var file_proto_transport_proto_depIdxs = []int32{
	4, // 0: transport.EventResponse.events:type_name -> transport.EventResponse.Event
	5, // 1: transport.AddRequest.events:type_name -> transport.AddRequest.EventData
	1, // 2: transport.EventStoreService.Add:input_type -> transport.AddRequest
	2, // 3: transport.EventStoreService.Get:input_type -> transport.GetRequest
	3, // 4: transport.EventStoreService.GetAll:input_type -> transport.GetAllRequest
	0, // 5: transport.EventStoreService.Add:output_type -> transport.EventResponse
	0, // 6: transport.EventStoreService.Get:output_type -> transport.EventResponse
	0, // 7: transport.EventStoreService.GetAll:output_type -> transport.EventResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_transport_proto_init() }
func file_proto_transport_proto_init() {
	if File_proto_transport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_transport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventResponse); i {
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
		file_proto_transport_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRequest); i {
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
		file_proto_transport_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_proto_transport_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAllRequest); i {
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
		file_proto_transport_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventResponse_Event); i {
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
		file_proto_transport_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddRequest_EventData); i {
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
			RawDescriptor: file_proto_transport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_transport_proto_goTypes,
		DependencyIndexes: file_proto_transport_proto_depIdxs,
		MessageInfos:      file_proto_transport_proto_msgTypes,
	}.Build()
	File_proto_transport_proto = out.File
	file_proto_transport_proto_rawDesc = nil
	file_proto_transport_proto_goTypes = nil
	file_proto_transport_proto_depIdxs = nil
}
