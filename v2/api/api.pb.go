// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

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

type AddEventsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddEventsRequest) Reset()         { *m = AddEventsRequest{} }
func (m *AddEventsRequest) String() string { return proto.CompactTextString(m) }
func (*AddEventsRequest) ProtoMessage()    {}
func (*AddEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *AddEventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddEventsRequest.Unmarshal(m, b)
}
func (m *AddEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddEventsRequest.Marshal(b, m, deterministic)
}
func (m *AddEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddEventsRequest.Merge(m, src)
}
func (m *AddEventsRequest) XXX_Size() int {
	return xxx_messageInfo_AddEventsRequest.Size(m)
}
func (m *AddEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddEventsRequest proto.InternalMessageInfo

type AddEventsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddEventsResponse) Reset()         { *m = AddEventsResponse{} }
func (m *AddEventsResponse) String() string { return proto.CompactTextString(m) }
func (*AddEventsResponse) ProtoMessage()    {}
func (*AddEventsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *AddEventsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddEventsResponse.Unmarshal(m, b)
}
func (m *AddEventsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddEventsResponse.Marshal(b, m, deterministic)
}
func (m *AddEventsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddEventsResponse.Merge(m, src)
}
func (m *AddEventsResponse) XXX_Size() int {
	return xxx_messageInfo_AddEventsResponse.Size(m)
}
func (m *AddEventsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddEventsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddEventsResponse proto.InternalMessageInfo

type GetEventsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetEventsRequest) Reset()         { *m = GetEventsRequest{} }
func (m *GetEventsRequest) String() string { return proto.CompactTextString(m) }
func (*GetEventsRequest) ProtoMessage()    {}
func (*GetEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *GetEventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetEventsRequest.Unmarshal(m, b)
}
func (m *GetEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetEventsRequest.Marshal(b, m, deterministic)
}
func (m *GetEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetEventsRequest.Merge(m, src)
}
func (m *GetEventsRequest) XXX_Size() int {
	return xxx_messageInfo_GetEventsRequest.Size(m)
}
func (m *GetEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetEventsRequest proto.InternalMessageInfo

type GetEventsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetEventsResponse) Reset()         { *m = GetEventsResponse{} }
func (m *GetEventsResponse) String() string { return proto.CompactTextString(m) }
func (*GetEventsResponse) ProtoMessage()    {}
func (*GetEventsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *GetEventsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetEventsResponse.Unmarshal(m, b)
}
func (m *GetEventsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetEventsResponse.Marshal(b, m, deterministic)
}
func (m *GetEventsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetEventsResponse.Merge(m, src)
}
func (m *GetEventsResponse) XXX_Size() int {
	return xxx_messageInfo_GetEventsResponse.Size(m)
}
func (m *GetEventsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetEventsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetEventsResponse proto.InternalMessageInfo

type LogEventsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEventsRequest) Reset()         { *m = LogEventsRequest{} }
func (m *LogEventsRequest) String() string { return proto.CompactTextString(m) }
func (*LogEventsRequest) ProtoMessage()    {}
func (*LogEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *LogEventsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEventsRequest.Unmarshal(m, b)
}
func (m *LogEventsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEventsRequest.Marshal(b, m, deterministic)
}
func (m *LogEventsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEventsRequest.Merge(m, src)
}
func (m *LogEventsRequest) XXX_Size() int {
	return xxx_messageInfo_LogEventsRequest.Size(m)
}
func (m *LogEventsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEventsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogEventsRequest proto.InternalMessageInfo

type LogEventsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogEventsResponse) Reset()         { *m = LogEventsResponse{} }
func (m *LogEventsResponse) String() string { return proto.CompactTextString(m) }
func (*LogEventsResponse) ProtoMessage()    {}
func (*LogEventsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *LogEventsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogEventsResponse.Unmarshal(m, b)
}
func (m *LogEventsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogEventsResponse.Marshal(b, m, deterministic)
}
func (m *LogEventsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogEventsResponse.Merge(m, src)
}
func (m *LogEventsResponse) XXX_Size() int {
	return xxx_messageInfo_LogEventsResponse.Size(m)
}
func (m *LogEventsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogEventsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogEventsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AddEventsRequest)(nil), "api.AddEventsRequest")
	proto.RegisterType((*AddEventsResponse)(nil), "api.AddEventsResponse")
	proto.RegisterType((*GetEventsRequest)(nil), "api.GetEventsRequest")
	proto.RegisterType((*GetEventsResponse)(nil), "api.GetEventsResponse")
	proto.RegisterType((*LogEventsRequest)(nil), "api.LogEventsRequest")
	proto.RegisterType((*LogEventsResponse)(nil), "api.LogEventsResponse")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 170 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0x12, 0xe2, 0x12, 0x70, 0x4c,
	0x49, 0x71, 0x2d, 0x4b, 0xcd, 0x2b, 0x29, 0x0e, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51, 0x12,
	0xe6, 0x12, 0x44, 0x12, 0x2b, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x05, 0x29, 0x74, 0x4f, 0x2d, 0xc1,
	0x50, 0x88, 0x24, 0x86, 0x50, 0xe8, 0x93, 0x9f, 0x8e, 0xa1, 0x10, 0x49, 0x0c, 0xa2, 0xd0, 0xe8,
	0x24, 0x23, 0x17, 0x6f, 0x70, 0x49, 0x51, 0x6a, 0x62, 0x6e, 0x70, 0x6a, 0x51, 0x59, 0x66, 0x72,
	0xaa, 0x90, 0x0d, 0x17, 0x27, 0xdc, 0x62, 0x21, 0x51, 0x3d, 0x90, 0x53, 0xd1, 0x1d, 0x27, 0x25,
	0x86, 0x2e, 0x0c, 0xb5, 0x96, 0x01, 0xa4, 0x1b, 0xee, 0x1a, 0xa8, 0x6e, 0x74, 0x17, 0x43, 0x75,
	0x63, 0x3a, 0x1a, 0xac, 0x1b, 0xee, 0x44, 0xa8, 0x6e, 0x74, 0x6f, 0x40, 0x75, 0x63, 0xf8, 0x44,
	0x89, 0xc1, 0x89, 0x35, 0x0a, 0x14, 0x9a, 0x49, 0x6c, 0xe0, 0x90, 0x35, 0x06, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x3f, 0xd3, 0x86, 0x00, 0x66, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamServiceClient interface {
	AddEvents(ctx context.Context, in *AddEventsRequest, opts ...grpc.CallOption) (*AddEventsResponse, error)
	GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error)
	LogEvents(ctx context.Context, in *LogEventsRequest, opts ...grpc.CallOption) (*LogEventsResponse, error)
}

type streamServiceClient struct {
	cc *grpc.ClientConn
}

func NewStreamServiceClient(cc *grpc.ClientConn) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) AddEvents(ctx context.Context, in *AddEventsRequest, opts ...grpc.CallOption) (*AddEventsResponse, error) {
	out := new(AddEventsResponse)
	err := c.cc.Invoke(ctx, "/api.StreamService/AddEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamServiceClient) GetEvents(ctx context.Context, in *GetEventsRequest, opts ...grpc.CallOption) (*GetEventsResponse, error) {
	out := new(GetEventsResponse)
	err := c.cc.Invoke(ctx, "/api.StreamService/GetEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *streamServiceClient) LogEvents(ctx context.Context, in *LogEventsRequest, opts ...grpc.CallOption) (*LogEventsResponse, error) {
	out := new(LogEventsResponse)
	err := c.cc.Invoke(ctx, "/api.StreamService/LogEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamServiceServer is the server API for StreamService service.
type StreamServiceServer interface {
	AddEvents(context.Context, *AddEventsRequest) (*AddEventsResponse, error)
	GetEvents(context.Context, *GetEventsRequest) (*GetEventsResponse, error)
	LogEvents(context.Context, *LogEventsRequest) (*LogEventsResponse, error)
}

// UnimplementedStreamServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStreamServiceServer struct {
}

func (*UnimplementedStreamServiceServer) AddEvents(ctx context.Context, req *AddEventsRequest) (*AddEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEvents not implemented")
}
func (*UnimplementedStreamServiceServer) GetEvents(ctx context.Context, req *GetEventsRequest) (*GetEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvents not implemented")
}
func (*UnimplementedStreamServiceServer) LogEvents(ctx context.Context, req *LogEventsRequest) (*LogEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogEvents not implemented")
}

func RegisterStreamServiceServer(s *grpc.Server, srv StreamServiceServer) {
	s.RegisterService(&_StreamService_serviceDesc, srv)
}

func _StreamService_AddEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServiceServer).AddEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StreamService/AddEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServiceServer).AddEvents(ctx, req.(*AddEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamService_GetEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServiceServer).GetEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StreamService/GetEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServiceServer).GetEvents(ctx, req.(*GetEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StreamService_LogEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamServiceServer).LogEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.StreamService/LogEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamServiceServer).LogEvents(ctx, req.(*LogEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StreamService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddEvents",
			Handler:    _StreamService_AddEvents_Handler,
		},
		{
			MethodName: "GetEvents",
			Handler:    _StreamService_GetEvents_Handler,
		},
		{
			MethodName: "LogEvents",
			Handler:    _StreamService_LogEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
