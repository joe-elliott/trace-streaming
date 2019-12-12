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
	ParentSpan
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
	TraceID       []byte      `protobuf:"bytes,1,opt,name=traceID,proto3" json:"traceID,omitempty"`
	SpanID        []byte      `protobuf:"bytes,2,opt,name=spanID,proto3" json:"spanID,omitempty"`
	ParentSpanID  []byte      `protobuf:"bytes,3,opt,name=parentSpanID,proto3" json:"parentSpanID,omitempty"`
	ProcessName   string      `protobuf:"bytes,4,opt,name=processName" json:"processName,omitempty"`
	OperationName string      `protobuf:"bytes,5,opt,name=operationName" json:"operationName,omitempty"`
	StartTime     int64       `protobuf:"varint,6,opt,name=startTime" json:"startTime,omitempty"`
	Duration      int64       `protobuf:"varint,7,opt,name=duration" json:"duration,omitempty"`
	Parent        *ParentSpan `protobuf:"bytes,8,opt,name=parent" json:"parent,omitempty"`
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

func (m *Span) GetSpanID() []byte {
	if m != nil {
		return m.SpanID
	}
	return nil
}

func (m *Span) GetParentSpanID() []byte {
	if m != nil {
		return m.ParentSpanID
	}
	return nil
}

func (m *Span) GetProcessName() string {
	if m != nil {
		return m.ProcessName
	}
	return ""
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

func (m *Span) GetParent() *ParentSpan {
	if m != nil {
		return m.Parent
	}
	return nil
}

type ParentSpan struct {
	ProcessName   string `protobuf:"bytes,1,opt,name=processName" json:"processName,omitempty"`
	OperationName string `protobuf:"bytes,2,opt,name=operationName" json:"operationName,omitempty"`
	StartTime     int64  `protobuf:"varint,3,opt,name=startTime" json:"startTime,omitempty"`
	Duration      int64  `protobuf:"varint,4,opt,name=duration" json:"duration,omitempty"`
}

func (m *ParentSpan) Reset()                    { *m = ParentSpan{} }
func (m *ParentSpan) String() string            { return proto.CompactTextString(m) }
func (*ParentSpan) ProtoMessage()               {}
func (*ParentSpan) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ParentSpan) GetProcessName() string {
	if m != nil {
		return m.ProcessName
	}
	return ""
}

func (m *ParentSpan) GetOperationName() string {
	if m != nil {
		return m.OperationName
	}
	return ""
}

func (m *ParentSpan) GetStartTime() int64 {
	if m != nil {
		return m.StartTime
	}
	return 0
}

func (m *ParentSpan) GetDuration() int64 {
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
	proto.RegisterType((*ParentSpan)(nil), "blergpb.ParentSpan")
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
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x6f, 0xd4, 0x30,
	0x10, 0xc5, 0xfb, 0xdd, 0x49, 0xf6, 0x32, 0x40, 0x65, 0x55, 0x1c, 0x56, 0x86, 0xc3, 0x4a, 0x48,
	0x2b, 0xb4, 0xdc, 0x90, 0x2a, 0xa4, 0xaa, 0x97, 0x1e, 0x40, 0xc5, 0xdb, 0x3f, 0xe0, 0xdd, 0x8c,
	0x20, 0x12, 0x89, 0x8d, 0xed, 0x1c, 0xe8, 0x5f, 0x00, 0xf1, 0x9b, 0x91, 0xed, 0x34, 0x69, 0x2a,
	0x58, 0x50, 0x6f, 0x99, 0xf7, 0x9e, 0x67, 0x5e, 0x9e, 0xc7, 0x90, 0xed, 0xbf, 0x92, 0xfd, 0xbc,
	0x31, 0x56, 0x7b, 0x8d, 0xf3, 0x58, 0x98, 0xbd, 0x20, 0x58, 0xee, 0xbc, 0x25, 0x55, 0x49, 0xfa,
	0xd6, 0x90, 0xf3, 0xb8, 0x01, 0xb4, 0xe9, 0x93, 0x8a, 0x0b, 0xe5, 0x0f, 0x5f, 0x76, 0xe5, 0x2d,
	0x71, 0xb6, 0x62, 0xeb, 0xa9, 0xfc, 0x03, 0x83, 0xaf, 0x60, 0xd9, 0xa1, 0x52, 0x79, 0xe2, 0xa3,
	0x28, 0x1d, 0x82, 0xe2, 0x16, 0xf2, 0x1b, 0xab, 0x0e, 0xd4, 0x4f, 0x99, 0x19, 0x65, 0x55, 0xe5,
	0x62, 0xe7, 0x6c, 0x7b, 0xba, 0x69, 0x0d, 0x6d, 0x06, 0x6e, 0x64, 0xab, 0xc2, 0x77, 0xc0, 0x0f,
	0x56, 0x3b, 0x47, 0xee, 0xda, 0xea, 0x03, 0x39, 0x77, 0xa1, 0x9b, 0xba, 0x50, 0xb6, 0x24, 0x17,
	0x07, 0x2e, 0xe4, 0x5f, 0x79, 0x71, 0x0e, 0xd9, 0xce, 0xa8, 0xfa, 0x91, 0xa3, 0xc5, 0x07, 0xc8,
	0xd3, 0x71, 0x67, 0x74, 0xed, 0x08, 0x39, 0xcc, 0x0b, 0xab, 0x8d, 0xa1, 0xa2, 0x4d, 0xe5, 0xae,
	0xc4, 0x97, 0x30, 0x75, 0x46, 0xd5, 0xc1, 0xd1, 0x78, 0x9d, 0x6d, 0x97, 0x7d, 0xe3, 0x70, 0x3e,
	0x71, 0xe2, 0xc7, 0x08, 0x26, 0xa1, 0x0e, 0x7d, 0x7c, 0x88, 0xe4, 0xea, 0x32, 0xf6, 0xc9, 0xe5,
	0x5d, 0x89, 0xa7, 0x30, 0x0b, 0xda, 0xab, 0xcb, 0xf8, 0x6b, 0xb9, 0x6c, 0x2b, 0x14, 0x90, 0x1b,
	0x65, 0xa9, 0xf6, 0xbb, 0xc4, 0x8e, 0x23, 0x3b, 0xc0, 0x70, 0x05, 0x99, 0x49, 0x09, 0x7c, 0x54,
	0x15, 0xf1, 0xc9, 0x8a, 0xad, 0x4f, 0xe4, 0x7d, 0x28, 0x5c, 0x98, 0x36, 0x64, 0x95, 0x2f, 0x75,
	0x1d, 0x35, 0xd3, 0xa8, 0x19, 0x82, 0xf8, 0x02, 0x4e, 0x9c, 0x57, 0xd6, 0xdf, 0x94, 0x15, 0xf1,
	0xd9, 0x8a, 0xad, 0xc7, 0xb2, 0x07, 0xf0, 0x0c, 0x16, 0x45, 0x93, 0xd4, 0x7c, 0x1e, 0xc9, 0xae,
	0xc6, 0xd7, 0x31, 0x5f, 0xaa, 0x3d, 0x5f, 0xc4, 0x7c, 0x9f, 0x76, 0x31, 0x5c, 0x77, 0x46, 0x65,
	0x2b, 0x11, 0xbf, 0x18, 0x40, 0x0f, 0x3f, 0x74, 0xcf, 0xfe, 0xc3, 0xfd, 0xe8, 0x9f, 0xee, 0xc7,
	0xc7, 0xdc, 0x4f, 0x86, 0xee, 0xb7, 0x3f, 0x19, 0x40, 0xb0, 0x92, 0x76, 0x01, 0xcf, 0x01, 0x3e,
	0x35, 0x64, 0xbf, 0x07, 0xc8, 0xe1, 0xb3, 0xe1, 0x8d, 0xa6, 0x45, 0x39, 0x7b, 0xfe, 0x00, 0x4d,
	0x7b, 0x22, 0x9e, 0xbc, 0x61, 0xf8, 0x1e, 0xb2, 0x78, 0x3c, 0xee, 0xbe, 0xc3, 0x5e, 0x79, 0xff,
	0x31, 0x1c, 0x69, 0xb0, 0x9f, 0xc5, 0xe7, 0xfa, 0xf6, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x43,
	0x01, 0x58, 0x3f, 0xbd, 0x03, 0x00, 0x00,
}
