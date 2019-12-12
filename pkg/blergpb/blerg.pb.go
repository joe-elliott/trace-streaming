// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blerg.proto

/*
Package blergpb is a generated protocol buffer package.

It is generated from these files:
	blerg.proto

It has these top-level messages:
	StreamRequest
	SpanResponse
	Span
*/
package blergpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StreamRequest struct {
	RequestedBatchSize int32 `protobuf:"varint,1,opt,name=requestedBatchSize" json:"requestedBatchSize,omitempty"`
	RequestedRate      int32 `protobuf:"varint,2,opt,name=requestedRate" json:"requestedRate,omitempty"`
}

func (m *StreamRequest) Reset()                    { *m = StreamRequest{} }
func (m *StreamRequest) String() string            { return proto.CompactTextString(m) }
func (*StreamRequest) ProtoMessage()               {}
func (*StreamRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StreamRequest) GetRequestedBatchSize() int32 {
	if m != nil {
		return m.RequestedBatchSize
	}
	return 0
}

func (m *StreamRequest) GetRequestedRate() int32 {
	if m != nil {
		return m.RequestedRate
	}
	return 0
}

type SpanResponse struct {
	Dropped int32   `protobuf:"varint,1,opt,name=dropped" json:"dropped,omitempty"`
	Spans   []*Span `protobuf:"bytes,2,rep,name=spans" json:"spans,omitempty"`
}

func (m *SpanResponse) Reset()                    { *m = SpanResponse{} }
func (m *SpanResponse) String() string            { return proto.CompactTextString(m) }
func (*SpanResponse) ProtoMessage()               {}
func (*SpanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SpanResponse) GetDropped() int32 {
	if m != nil {
		return m.Dropped
	}
	return 0
}

func (m *SpanResponse) GetSpans() []*Span {
	if m != nil {
		return m.Spans
	}
	return nil
}

type Span struct {
	OperationName string `protobuf:"bytes,1,opt,name=operationName" json:"operationName,omitempty"`
	StartTime     int32  `protobuf:"varint,2,opt,name=startTime" json:"startTime,omitempty"`
	Duration      int32  `protobuf:"varint,3,opt,name=duration" json:"duration,omitempty"`
	Parent        *Span  `protobuf:"bytes,4,opt,name=parent" json:"parent,omitempty"`
}

func (m *Span) Reset()                    { *m = Span{} }
func (m *Span) String() string            { return proto.CompactTextString(m) }
func (*Span) ProtoMessage()               {}
func (*Span) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Span) GetOperationName() string {
	if m != nil {
		return m.OperationName
	}
	return ""
}

func (m *Span) GetStartTime() int32 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *Span) GetDuration() int32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *Span) GetParent() *Span {
	if m != nil {
		return m.Parent
	}
	return nil
}

func init() {
	proto.RegisterType((*StreamRequest)(nil), "blergpb.StreamRequest")
	proto.RegisterType((*SpanResponse)(nil), "blergpb.SpanResponse")
	proto.RegisterType((*Span)(nil), "blergpb.Span")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SpanStream service

type SpanStreamClient interface {
	Tail(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (SpanStream_TailClient, error)
}

type spanStreamClient struct {
	cc *grpc.ClientConn
}

func NewSpanStreamClient(cc *grpc.ClientConn) SpanStreamClient {
	return &spanStreamClient{cc}
}

func (c *spanStreamClient) Tail(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (SpanStream_TailClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SpanStream_serviceDesc.Streams[0], c.cc, "/blergpb.SpanStream/Tail", opts...)
	if err != nil {
		return nil, err
	}
	x := &spanStreamTailClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SpanStream_TailClient interface {
	Recv() (*SpanResponse, error)
	grpc.ClientStream
}

type spanStreamTailClient struct {
	grpc.ClientStream
}

func (x *spanStreamTailClient) Recv() (*SpanResponse, error) {
	m := new(SpanResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SpanStream service

type SpanStreamServer interface {
	Tail(*StreamRequest, SpanStream_TailServer) error
}

func RegisterSpanStreamServer(s *grpc.Server, srv SpanStreamServer) {
	s.RegisterService(&_SpanStream_serviceDesc, srv)
}

func _SpanStream_Tail_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SpanStreamServer).Tail(m, &spanStreamTailServer{stream})
}

type SpanStream_TailServer interface {
	Send(*SpanResponse) error
	grpc.ServerStream
}

type spanStreamTailServer struct {
	grpc.ServerStream
}

func (x *spanStreamTailServer) Send(m *SpanResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _SpanStream_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blergpb.SpanStream",
	HandlerType: (*SpanStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Tail",
			Handler:       _SpanStream_Tail_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "blerg.proto",
}

func init() { proto.RegisterFile("blerg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x3f, 0x4f, 0xc3, 0x30,
	0x10, 0xc5, 0x49, 0xff, 0xd2, 0x2b, 0x59, 0x4e, 0x02, 0x59, 0x15, 0x43, 0x14, 0x40, 0xea, 0x14,
	0xa1, 0x32, 0xb1, 0xb2, 0x30, 0xc1, 0xe0, 0xf4, 0x0b, 0x5c, 0xc8, 0x09, 0x22, 0x35, 0xb6, 0xb1,
	0xdd, 0x85, 0xef, 0xc0, 0x77, 0x46, 0xb1, 0xd3, 0x96, 0x08, 0x36, 0xff, 0xde, 0x3b, 0x9f, 0x9f,
	0x9e, 0x61, 0x59, 0xed, 0xd8, 0xbe, 0x17, 0xc6, 0x6a, 0xaf, 0x71, 0x1e, 0xc0, 0x54, 0x39, 0x43,
	0x5a, 0x7a, 0xcb, 0xd4, 0x4a, 0xfe, 0xdc, 0xb3, 0xf3, 0x58, 0x00, 0xda, 0x78, 0xe4, 0xfa, 0x89,
	0xfc, 0xdb, 0x47, 0xd9, 0x7c, 0xb1, 0x48, 0xb2, 0x64, 0x3d, 0x95, 0xff, 0x38, 0x78, 0x0b, 0xe9,
	0x51, 0x95, 0xe4, 0x59, 0x8c, 0xc2, 0xe8, 0x50, 0xcc, 0x5f, 0xe0, 0xa2, 0x34, 0xa4, 0x24, 0x3b,
	0xa3, 0x95, 0x63, 0x14, 0x30, 0xaf, 0xad, 0x36, 0x86, 0xeb, 0x7e, 0xf5, 0x01, 0xf1, 0x06, 0xa6,
	0xce, 0x90, 0x72, 0x62, 0x94, 0x8d, 0xd7, 0xcb, 0x4d, 0x5a, 0xf4, 0x49, 0x8b, 0x70, 0x3f, 0x7a,
	0xf9, 0x77, 0x02, 0x93, 0x8e, 0xbb, 0xd7, 0xb5, 0x61, 0x4b, 0xbe, 0xd1, 0xea, 0x95, 0xda, 0x18,
	0x74, 0x21, 0x87, 0x22, 0x5e, 0xc3, 0xc2, 0x79, 0xb2, 0x7e, 0xdb, 0xb4, 0x87, 0x7c, 0x27, 0x01,
	0x57, 0x70, 0x5e, 0xef, 0xe3, 0xb4, 0x18, 0x07, 0xf3, 0xc8, 0x78, 0x07, 0x33, 0x43, 0x96, 0x95,
	0x17, 0x93, 0x2c, 0xf9, 0x1b, 0xa7, 0x37, 0x37, 0xcf, 0x00, 0x1d, 0xc7, 0x26, 0xf1, 0x11, 0x26,
	0x5b, 0x6a, 0x76, 0x78, 0x75, 0x1a, 0xfe, 0x5d, 0xf1, 0xea, 0x72, 0xb8, 0xa4, 0xef, 0x24, 0x3f,
	0xbb, 0x4f, 0xaa, 0x59, 0xf8, 0x9e, 0x87, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x14, 0x59, 0x5d,
	0x61, 0xad, 0x01, 0x00, 0x00,
}