// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blerg.proto

/*
Package blergpb is a generated protocol buffer package.

It is generated from these files:
	blerg.proto

It has these top-level messages:
	StreamRequest
	TraceRequest
	SpanRequest
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

type TraceRequest struct {
	Params                   *StreamRequest `protobuf:"bytes,1,opt,name=params" json:"params,omitempty"`
	CrossesProcessBoundaries bool           `protobuf:"varint,2,opt,name=crossesProcessBoundaries" json:"crossesProcessBoundaries,omitempty"`
}

func (m *TraceRequest) Reset()                    { *m = TraceRequest{} }
func (m *TraceRequest) String() string            { return proto.CompactTextString(m) }
func (*TraceRequest) ProtoMessage()               {}
func (*TraceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TraceRequest) GetParams() *StreamRequest {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *TraceRequest) GetCrossesProcessBoundaries() bool {
	if m != nil {
		return m.CrossesProcessBoundaries
	}
	return false
}

type SpanRequest struct {
	Params *StreamRequest `protobuf:"bytes,1,opt,name=params" json:"params,omitempty"`
}

func (m *SpanRequest) Reset()                    { *m = SpanRequest{} }
func (m *SpanRequest) String() string            { return proto.CompactTextString(m) }
func (*SpanRequest) ProtoMessage()               {}
func (*SpanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *SpanRequest) GetParams() *StreamRequest {
	if m != nil {
		return m.Params
	}
	return nil
}

type SpanResponse struct {
	Dropped int32   `protobuf:"varint,1,opt,name=dropped" json:"dropped,omitempty"`
	Spans   []*Span `protobuf:"bytes,2,rep,name=spans" json:"spans,omitempty"`
}

func (m *SpanResponse) Reset()                    { *m = SpanResponse{} }
func (m *SpanResponse) String() string            { return proto.CompactTextString(m) }
func (*SpanResponse) ProtoMessage()               {}
func (*SpanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
	TraceID       []byte `protobuf:"bytes,1,opt,name=traceID,proto3" json:"traceID,omitempty"`
	OperationName string `protobuf:"bytes,2,opt,name=operationName" json:"operationName,omitempty"`
	StartTime     int64  `protobuf:"varint,3,opt,name=startTime" json:"startTime,omitempty"`
	Duration      int64  `protobuf:"varint,4,opt,name=duration" json:"duration,omitempty"`
}

func (m *Span) Reset()                    { *m = Span{} }
func (m *Span) String() string            { return proto.CompactTextString(m) }
func (*Span) ProtoMessage()               {}
func (*Span) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Span) GetTraceID() []byte {
	if m != nil {
		return m.TraceID
	}
	return nil
}

func (m *Span) GetOperationName() string {
	if m != nil {
		return m.OperationName
	}
	return ""
}

func (m *Span) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *Span) GetDuration() int64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func init() {
	proto.RegisterType((*StreamRequest)(nil), "blergpb.StreamRequest")
	proto.RegisterType((*TraceRequest)(nil), "blergpb.TraceRequest")
	proto.RegisterType((*SpanRequest)(nil), "blergpb.SpanRequest")
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
	QuerySpans(ctx context.Context, in *SpanRequest, opts ...grpc.CallOption) (SpanStream_QuerySpansClient, error)
	QueryTraces(ctx context.Context, in *TraceRequest, opts ...grpc.CallOption) (SpanStream_QueryTracesClient, error)
}

type spanStreamClient struct {
	cc *grpc.ClientConn
}

func NewSpanStreamClient(cc *grpc.ClientConn) SpanStreamClient {
	return &spanStreamClient{cc}
}

func (c *spanStreamClient) QuerySpans(ctx context.Context, in *SpanRequest, opts ...grpc.CallOption) (SpanStream_QuerySpansClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SpanStream_serviceDesc.Streams[0], c.cc, "/blergpb.SpanStream/QuerySpans", opts...)
	if err != nil {
		return nil, err
	}
	x := &spanStreamQuerySpansClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SpanStream_QuerySpansClient interface {
	Recv() (*SpanResponse, error)
	grpc.ClientStream
}

type spanStreamQuerySpansClient struct {
	grpc.ClientStream
}

func (x *spanStreamQuerySpansClient) Recv() (*SpanResponse, error) {
	m := new(SpanResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *spanStreamClient) QueryTraces(ctx context.Context, in *TraceRequest, opts ...grpc.CallOption) (SpanStream_QueryTracesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SpanStream_serviceDesc.Streams[1], c.cc, "/blergpb.SpanStream/QueryTraces", opts...)
	if err != nil {
		return nil, err
	}
	x := &spanStreamQueryTracesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SpanStream_QueryTracesClient interface {
	Recv() (*SpanResponse, error)
	grpc.ClientStream
}

type spanStreamQueryTracesClient struct {
	grpc.ClientStream
}

func (x *spanStreamQueryTracesClient) Recv() (*SpanResponse, error) {
	m := new(SpanResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SpanStream service

type SpanStreamServer interface {
	QuerySpans(*SpanRequest, SpanStream_QuerySpansServer) error
	QueryTraces(*TraceRequest, SpanStream_QueryTracesServer) error
}

func RegisterSpanStreamServer(s *grpc.Server, srv SpanStreamServer) {
	s.RegisterService(&_SpanStream_serviceDesc, srv)
}

func _SpanStream_QuerySpans_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SpanRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SpanStreamServer).QuerySpans(m, &spanStreamQuerySpansServer{stream})
}

type SpanStream_QuerySpansServer interface {
	Send(*SpanResponse) error
	grpc.ServerStream
}

type spanStreamQuerySpansServer struct {
	grpc.ServerStream
}

func (x *spanStreamQuerySpansServer) Send(m *SpanResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SpanStream_QueryTraces_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TraceRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SpanStreamServer).QueryTraces(m, &spanStreamQueryTracesServer{stream})
}

type SpanStream_QueryTracesServer interface {
	Send(*SpanResponse) error
	grpc.ServerStream
}

type spanStreamQueryTracesServer struct {
	grpc.ServerStream
}

func (x *spanStreamQueryTracesServer) Send(m *SpanResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _SpanStream_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blergpb.SpanStream",
	HandlerType: (*SpanStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "QuerySpans",
			Handler:       _SpanStream_QuerySpans_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "QueryTraces",
			Handler:       _SpanStream_QueryTraces_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "blerg.proto",
}

func init() { proto.RegisterFile("blerg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x3d, 0x4f, 0xc3, 0x30,
	0x10, 0x25, 0xf4, 0xfb, 0xd2, 0x2e, 0x16, 0xa0, 0xa8, 0x62, 0xa8, 0x02, 0x43, 0xa7, 0x08, 0x95,
	0x0d, 0xa9, 0x42, 0xaa, 0x58, 0x18, 0x40, 0xe0, 0xf6, 0x0f, 0xb8, 0xc9, 0x09, 0x22, 0x91, 0xd8,
	0xf8, 0x9c, 0x81, 0x4e, 0xfc, 0x00, 0x7e, 0x34, 0xb2, 0xdd, 0x36, 0x04, 0x01, 0x03, 0x5b, 0xde,
	0xbb, 0xbb, 0x77, 0x2f, 0xef, 0x0c, 0xe1, 0xfa, 0x05, 0xf5, 0x53, 0xa2, 0xb4, 0x34, 0x92, 0xf5,
	0x1c, 0x50, 0xeb, 0x18, 0x61, 0xb4, 0x34, 0x1a, 0x45, 0xc1, 0xf1, 0xb5, 0x42, 0x32, 0x2c, 0x01,
	0xa6, 0xfd, 0x27, 0x66, 0x0b, 0x61, 0xd2, 0xe7, 0x65, 0xbe, 0xc1, 0x28, 0x98, 0x04, 0xd3, 0x0e,
	0xff, 0xa1, 0xc2, 0xce, 0x61, 0xb4, 0x67, 0xb9, 0x30, 0x18, 0x1d, 0xba, 0xd6, 0x26, 0x19, 0x6f,
	0x60, 0xb8, 0xd2, 0x22, 0xc5, 0x7a, 0x4b, 0x57, 0x09, 0x2d, 0x0a, 0x72, 0xca, 0xe1, 0xec, 0x24,
	0xd9, 0x1a, 0x4a, 0x1a, 0x6e, 0xf8, 0xb6, 0x8b, 0x5d, 0x41, 0x94, 0x6a, 0x49, 0x84, 0xf4, 0xa0,
	0x65, 0x8a, 0x44, 0x0b, 0x59, 0x95, 0x99, 0xd0, 0x39, 0x92, 0x5b, 0xd8, 0xe7, 0xbf, 0xd6, 0xe3,
	0x39, 0x84, 0x4b, 0x25, 0xca, 0x7f, 0xae, 0x8e, 0xef, 0x60, 0xe8, 0xc7, 0x49, 0xc9, 0x92, 0x90,
	0x45, 0xd0, 0xcb, 0xb4, 0x54, 0x0a, 0xb3, 0x6d, 0x2a, 0x3b, 0xc8, 0xce, 0xa0, 0x43, 0x4a, 0x94,
	0xd6, 0x51, 0x6b, 0x1a, 0xce, 0x46, 0xb5, 0xb0, 0x9d, 0xf7, 0xb5, 0xf8, 0x3d, 0x80, 0xb6, 0xc5,
	0x56, 0xc7, 0xd8, 0x48, 0x6e, 0x6f, 0x9c, 0xce, 0x90, 0xef, 0xa0, 0x8d, 0x54, 0x2a, 0xd4, 0xc2,
	0xe4, 0xb2, 0xbc, 0x17, 0x85, 0x8f, 0x74, 0xc0, 0x9b, 0x24, 0x3b, 0x85, 0x01, 0x19, 0xa1, 0xcd,
	0x2a, 0x2f, 0x30, 0x6a, 0x4d, 0x82, 0x69, 0x8b, 0xd7, 0x04, 0x1b, 0x43, 0x3f, 0xab, 0x7c, 0x77,
	0xd4, 0x76, 0xc5, 0x3d, 0x9e, 0x7d, 0x04, 0x00, 0xd6, 0x82, 0xff, 0x5f, 0x36, 0x07, 0x78, 0xac,
	0x50, 0xbf, 0x59, 0x8a, 0xd8, 0x51, 0xd3, 0xb5, 0x0f, 0x63, 0x7c, 0xfc, 0x8d, 0xf5, 0x59, 0xc4,
	0x07, 0x17, 0x01, 0xbb, 0x86, 0xd0, 0x8d, 0xbb, 0xfb, 0x12, 0xab, 0x3b, 0xbf, 0x1e, 0xfc, 0x0f,
	0x81, 0x75, 0xd7, 0x3d, 0xc9, 0xcb, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x87, 0x94, 0x5f, 0xf2,
	0xa1, 0x02, 0x00, 0x00,
}
