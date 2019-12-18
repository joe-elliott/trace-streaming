package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

type GRPCConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

type grpcServer struct {
	s      StreamProcessor
	cfg    GRPCConfig
	server *grpc.Server
}

func NewGRPC(s StreamProcessor, cfg GRPCConfig) StreamServer {
	return &grpcServer{
		s:   s,
		cfg: cfg,
	}
}

func (g *grpcServer) Do() error {
	if !g.cfg.Enabled {
		return nil
	}

	grpcEndpoint := fmt.Sprintf(":%d", g.cfg.Port)
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return fmt.Errorf("grpc failed to listen %v", err)
	}

	g.server = grpc.NewServer()
	streampb.RegisterSpanStreamServer(g.server, g)
	go func() {
		err := g.server.Serve(lis)
		if err != nil {
			log.Fatal("Failed to start GRPC Server", err)
		}
	}()

	return nil
}

func (g *grpcServer) Shutdown() {
	if g.server != nil {
		g.server.Stop()
	}
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
