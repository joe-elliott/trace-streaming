package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/util"
)

type grpcServer struct {
	s StreamProcessor
}

func DoGRPC(s StreamProcessor) error {
	g := &grpcServer{
		s: s,
	}

	// GRPC
	grpcEndpoint := fmt.Sprintf(":%d", util.DefaultGRPCPort)
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	server := grpc.NewServer()
	streampb.RegisterSpanStreamServer(server, g)
	go func() {
		err := server.Serve(lis)
		if err != nil {
			log.Fatal("Failed to start GRPC Server", err)
		}
	}()

	return nil
}

func (g *grpcServer) QuerySpans(req *streampb.SpanRequest, stream streampb.SpanStream_QuerySpansServer) error {
	tailer := streamer.NewSpans(req, stream)

	g.s.AddSpanStreamer(tailer)

	return nil
}

func (g *grpcServer) QueryTraces(req *streampb.TraceRequest, stream streampb.SpanStream_QueryTracesServer) error {
	tailer := streamer.NewTraces(req, stream)

	g.s.AddTraceStreamer(tailer)

	return nil
}
