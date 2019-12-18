package streamprocessor

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/open-telemetry/opentelemetry-collector/consumer"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/server"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

type streamProcessor struct {
	nextConsumer consumer.TraceConsumer
	config       Config

	spanStreamers  []*streamer.Spans
	traceStreamers []*streamer.Traces

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

	server.DoGRPC(sp, config.GRPC)
	server.DoWebsocket(sp, config.Websocket)

	go sp.pollBatches(5 * time.Second)

	return sp, nil
}

func (sp *streamProcessor) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {
	blergSpans := make([]*streampb.Span, len(td.Spans))
	i := 0

	for _, span := range td.Spans {
		if !isSpanValid(span) {
			continue
		}

		blergSpan := spanToSpan(span, td.Node)
		blergSpans[i] = blergSpan
		i++
	}

	if i > 0 {
		blergSpans = blergSpans[:i-1]

		for _, s := range sp.spanStreamers {
			s.ProcessBatch(blergSpans)
		}

		sp.traceBatcher.addBatch(blergSpans)
	}

	return sp.nextConsumer.ConsumeTraceData(ctx, td)
}

func (sp *streamProcessor) GetCapabilities() processor.Capabilities {
	return processor.Capabilities{MutatesConsumedData: false}
}

func (sp *streamProcessor) Shutdown() error {
	for _, s := range sp.spanStreamers {
		s.Shutdown()
	}

	for _, s := range sp.traceStreamers {
		s.Shutdown()
	}

	return nil
}

func (sp *streamProcessor) AddSpanStreamer(s *streamer.Spans) {
	sp.spanStreamers = append(sp.spanStreamers, s)

	s.Do()
}

func (sp *streamProcessor) AddTraceStreamer(s *streamer.Traces) {
	sp.traceStreamers = append(sp.traceStreamers, s)

	s.Do()
}

func (sp *streamProcessor) pollBatches(pollTime time.Duration) {
	ticker := time.NewTicker(pollTime)

	for {
		completed := sp.traceBatcher.completeBatches()

		for _, batch := range completed {
			tree := buildSpanTree(batch)

			for _, t := range sp.traceStreamers {
				t.ProcessBatch(tree)
			}
		}

		<-ticker.C
	}
}

func isSpanValid(span *tracepb.Span) bool {
	return span != nil && len(span.TraceId) > 0 && len(span.SpanId) > 0
}

func spanToSpan(in *tracepb.Span, node *commonpb.Node) *streampb.Span {
	name := "unknown"
	processName := "unknown"
	startTime := int64(0)
	duration := int32(0)

	if node != nil && node.ServiceInfo != nil {
		processName = node.ServiceInfo.Name
	}

	if in.Name != nil {
		name = in.Name.Value
	}

	if in.StartTime != nil {
		startTime = in.StartTime.Seconds

		if in.EndTime != nil {
			duration = int32((in.EndTime.Nanos - in.StartTime.Nanos) / 1000000)
		}
	}

	return &streampb.Span{
		TraceID:       in.TraceId,
		SpanID:        in.SpanId,
		ParentSpanID:  in.ParentSpanId,
		ProcessName:   processName,
		OperationName: name,
		StartTime:     startTime,
		Duration:      duration,
	}
}

func buildSpanTree(trace []*streampb.Span) []*streampb.Span {
	tree := make([]*streampb.Span, 0)

	// O(n^2)! yay!
	for _, child := range trace {

		found := false
		for _, parent := range trace {

			if bytes.Equal(child.ParentSpanID, parent.SpanID) {
				found = true

				child.Parent = &streampb.ParentSpan{
					OperationName: parent.OperationName,
					ProcessName:   parent.ProcessName,
					StartTime:     parent.StartTime,
					Duration:      parent.Duration,
				}
			}
		}

		if !found && len(child.ParentSpanID) > 0 {
			log.Printf("Unable to find parent id %v. Dropping.", child.ParentSpanID)
			continue
		}

		tree = append(tree, child)
	}

	return tree
}
