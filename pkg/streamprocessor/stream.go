package streamprocessor

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/open-telemetry/opentelemetry-collector/consumer"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"

	"github.com/joe-elliott/blerg/pkg/blergpb"
	"github.com/joe-elliott/blerg/pkg/streamer"
	"github.com/joe-elliott/blerg/pkg/util"
)

type streamProcessor struct {
	nextConsumer  consumer.TraceConsumer
	config        Config
	spanStreamers []*streamer.Spans
}

// NewTraceProcessor returns the span processor.
func NewTraceProcessor(nextConsumer consumer.TraceConsumer, config Config) (processor.TraceProcessor, error) {
	if nextConsumer == nil {
		return nil, oterr.ErrNilNextConsumer
	}

	sp := &streamProcessor{
		nextConsumer: nextConsumer,
		config:       config,
	}

	port := util.DefaultPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	server := grpc.NewServer()

	blergpb.RegisterSpanStreamServer(server, sp)

	go func() {
		go server.Serve(lis)
	}()

	return sp, nil
}

func (sp *streamProcessor) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {

	for _, s := range sp.spanStreamers {
		s.ProcessBatch(td.Spans)
	}

	return sp.nextConsumer.ConsumeTraceData(ctx, td)
}

func (sp *streamProcessor) GetCapabilities() processor.Capabilities {
	return processor.Capabilities{MutatesConsumedData: false}
}

func (sp *streamProcessor) Tail(req *blergpb.StreamRequest, stream blergpb.SpanStream_TailServer) error {
	tailer := streamer.NewSpans(req, stream)
	sp.spanStreamers = append(sp.spanStreamers, tailer)

	return tailer.Do()
}
