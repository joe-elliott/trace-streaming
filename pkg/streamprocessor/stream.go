package streamprocessor

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/open-telemetry/opentelemetry-collector/consumer"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"

	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"

	"github.com/joe-elliott/blerg/pkg/blergpb"
	"github.com/joe-elliott/blerg/pkg/streamer"
	"github.com/joe-elliott/blerg/pkg/util"
)

type streamProcessor struct {
	nextConsumer consumer.TraceConsumer
	config       Config

	spanStreamers  []*streamer.Spans
	traceStreamers []*streamer.Spans

	traceBatcher *batcher
}

// NewTraceProcessor returns the span processor.
func NewTraceProcessor(nextConsumer consumer.TraceConsumer, config Config) (processor.TraceProcessor, error) {
	if nextConsumer == nil {
		return nil, oterr.ErrNilNextConsumer
	}

	sp := &streamProcessor{
		nextConsumer: nextConsumer,
		config:       config,
		traceBatcher: newBatcher(),
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

	go sp.pollBatches(5 * time.Second)

	return sp, nil
}

func (sp *streamProcessor) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {
	blergSpans := make([]*blergpb.Span, len(td.Spans))

	for i, span := range td.Spans {
		blergSpan := spanToSpan(span)
		blergSpans[i] = blergSpan
	}

	for _, s := range sp.spanStreamers {
		s.ProcessBatch(blergSpans)
	}

	sp.traceBatcher.addBatch(blergSpans)

	return sp.nextConsumer.ConsumeTraceData(ctx, td)
}

func (sp *streamProcessor) GetCapabilities() processor.Capabilities {
	return processor.Capabilities{MutatesConsumedData: false}
}

func (sp *streamProcessor) QuerySpans(req *blergpb.StreamRequest, stream blergpb.SpanStream_QuerySpansServer) error {
	tailer := streamer.NewSpans(req, stream)
	sp.spanStreamers = append(sp.spanStreamers, tailer)

	return tailer.Do()
}

func (sp *streamProcessor) QueryTraces(req *blergpb.StreamRequest, stream blergpb.SpanStream_QueryTracesServer) error {
	tailer := streamer.NewSpans(req, stream)
	sp.traceStreamers = append(sp.traceStreamers, tailer)

	return tailer.Do()
}

func (sp *streamProcessor) pollBatches(pollTime time.Duration) {
	ticker := time.NewTicker(pollTime)

	for {
		completed := sp.traceBatcher.completeBatches()

		for _, batch := range completed {
			for _, t := range sp.traceStreamers {
				t.ProcessBatch(batch)
			}
		}

		<-ticker.C
	}
}

func spanToSpan(in *tracepb.Span) *blergpb.Span {
	return &blergpb.Span{
		TraceID:       in.TraceId,
		OperationName: in.Name.Value,
		StartTime:     in.StartTime.Seconds,
		Duration:      in.EndTime.Seconds - in.StartTime.Seconds,
	}
}
