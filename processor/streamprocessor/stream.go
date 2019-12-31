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

	spanStreamers  []streamer.Streamer
	traceStreamers []streamer.Streamer

	traceBatcher *batcher

	servers []server.StreamServer
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

	sp.servers = append(sp.servers, server.NewGRPC(sp, config.GRPC))
	sp.servers = append(sp.servers, server.NewWebsocket(sp, config.Websocket))

	for _, srv := range sp.servers {
		err := srv.Do()

		if err != nil {
			return nil, err
		}
	}

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
		blergSpans = blergSpans[:i]
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

	for _, srv := range sp.servers {
		srv.Shutdown()
	}

	return nil
}

func (sp *streamProcessor) AddSpanStreamer(s streamer.Streamer) {
	sp.spanStreamers = append(sp.spanStreamers, s)

	s.Do()
}

func (sp *streamProcessor) AddTraceStreamer(s streamer.Streamer) {
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

	var status *streampb.Status
	if in.Status != nil {
		status = &streampb.Status{
			Code:    streampb.Status_StatusCode(in.Status.Code),
			Message: in.Status.Message,
		}
	}

	return &streampb.Span{
		Name:         name,
		TraceID:      in.TraceId,
		SpanID:       in.SpanId,
		ParentSpanID: in.ParentSpanId,
		Process: &streampb.Process{
			Name: processName,
		},
		Status:      status,
		Events:      nil, //todo: support events
		Attributes:  attributesToKVP(in.Attributes),
		StartTime:   startTime,
		Duration:    duration,
		ParentIndex: -1,
	}
}

func attributesToKVP(atts *tracepb.Span_Attributes) map[string]*streampb.KeyValuePair {
	if atts == nil {
		return nil
	}

	ret := make(map[string]*streampb.KeyValuePair)

	for k, v := range atts.AttributeMap {
		kvp := &streampb.KeyValuePair{
			Key: k,
		}

		switch val := v.Value.(type) {
		case *tracepb.AttributeValue_StringValue:
			kvp.StringValue = val.StringValue.Value
			kvp.Type = streampb.KeyValuePair_STRING
		case *tracepb.AttributeValue_IntValue:
			kvp.IntValue = val.IntValue
			kvp.Type = streampb.KeyValuePair_INT
		case *tracepb.AttributeValue_BoolValue:
			kvp.BoolValue = val.BoolValue
			kvp.Type = streampb.KeyValuePair_BOOL
		case *tracepb.AttributeValue_DoubleValue:
			kvp.DoubleValue = val.DoubleValue
			kvp.Type = streampb.KeyValuePair_DOUBLE
		}

		ret[k] = kvp
	}

	return ret
}

// todo: make root span at position 0
func buildSpanTree(trace []*streampb.Span) []*streampb.Span {
	tree := make([]*streampb.Span, 0)

	// O(n^2)! yay!
	for _, child := range trace {

		found := false
		for i, parent := range trace {

			if bytes.Equal(child.ParentSpanID, parent.SpanID) {
				found = true

				child.ParentIndex = int32(i)
			}
		}

		// todo: remove this kludge
		if !found && len(child.ParentSpanID) > 0 {
			log.Printf("Unable to find parent id %v. Dropping.", child.ParentSpanID)
			continue
		}

		tree = append(tree, child)
	}

	return tree
}
