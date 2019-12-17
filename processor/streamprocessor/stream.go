package streamprocessor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/open-telemetry/opentelemetry-collector/consumer"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"

	commonpb "github.com/census-instrumentation/opencensus-proto/gen-go/agent/common/v1"
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"

	"github.com/joe-elliott/blerg/processor/streamprocessor/streampb"
	"github.com/joe-elliott/blerg/processor/streamprocessor/streamer"
	"github.com/joe-elliott/blerg/processor/streamprocessor/util"

	"github.com/gorilla/websocket"
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

	// GRPC
	grpcEndpoint := fmt.Sprintf(":%d", util.DefaultGRPCPort)
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	server := grpc.NewServer()
	streampb.RegisterSpanStreamServer(server, sp)
	go func() {
		err := server.Serve(lis)
		if err != nil {
			log.Fatal("Failed to start GRPC Server", err)
		}
	}()

	go sp.pollBatches(5 * time.Second)

	sp.startWebsocket()

	return sp, nil
}

func (sp *streamProcessor) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {
	blergSpans := make([]*streampb.Span, len(td.Spans))

	for i, span := range td.Spans {
		blergSpan := spanToSpan(span, td.Node)
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

func (sp *streamProcessor) Shutdown() error {
	for _, s := range sp.spanStreamers {
		s.Shutdown()
	}

	for _, s := range sp.traceStreamers {
		s.Shutdown()
	}

	return nil
}

func (sp *streamProcessor) QuerySpans(req *streampb.SpanRequest, stream streampb.SpanStream_QuerySpansServer) error {
	tailer := streamer.NewSpans(req, stream)
	sp.spanStreamers = append(sp.spanStreamers, tailer)

	return tailer.Do()
}

func (sp *streamProcessor) QueryTraces(req *streampb.TraceRequest, stream streampb.SpanStream_QueryTracesServer) error {
	tailer := streamer.NewTraces(req, stream)
	sp.traceStreamers = append(sp.traceStreamers, tailer)

	return tailer.Do()
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

// websocket crap
type socketSender struct {
	ws *websocket.Conn
}

func (sp *streamProcessor) startWebsocket() {
	http.HandleFunc("/v1/stream/traces", sp.streamTraces)
	http.HandleFunc("/v1/stream/spans", sp.streamSpans)
	go http.ListenAndServe(fmt.Sprintf(":%d", util.DefaultHTTPPort), nil)
}

func (sp *streamProcessor) streamTraces(w http.ResponseWriter, r *http.Request) {
	s := setupWebsocket(w, r)

	query := r.URL.Query()

	tailer := streamer.NewTraces(&streampb.TraceRequest{
		Params:      getStreamRequest(query),
		ProcessName: getQueryParam(query, "processName"),
	}, s)
	sp.traceStreamers = append(sp.traceStreamers, tailer)

	tailer.Do()
}

func (sp *streamProcessor) streamSpans(w http.ResponseWriter, r *http.Request) {
	s := setupWebsocket(w, r)

	query := r.URL.Query()

	tailer := streamer.NewSpans(&streampb.SpanRequest{
		Params:        getStreamRequest(query),
		ProcessName:   getQueryParam(query, "processName"),
		OperationName: getQueryParam(query, "operationName"),
		MinDuration:   int32(getQueryParamInt(query, "minDuration")),
	}, s)
	sp.spanStreamers = append(sp.spanStreamers, tailer)

	tailer.Do()
}

func (s *socketSender) Send(span *streampb.SpanResponse) error {
	return s.ws.WriteJSON(span)
}

// utility
func getQueryParam(v url.Values, name string) string {
	value, ok := v[name]

	if ok && len(value) > 0 {
		return value[0]
	}

	return ""
}

func getQueryParamInt(v url.Values, name string) int {
	value, ok := v[name]

	ret := 0

	if ok && len(value) > 0 {
		ret, _ = strconv.Atoi(value[0])
	}

	return ret
}

func getStreamRequest(v url.Values) *streampb.StreamRequest {
	rate := getQueryParamInt(v, "rate")

	return &streampb.StreamRequest{
		RequestedRate: int32(rate),
	}
}

func setupWebsocket(w http.ResponseWriter, r *http.Request) *socketSender {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// helpful log statement to show connections
	log.Println("Client Connected")

	return &socketSender{
		ws: ws,
	}
}

func spanToSpan(in *tracepb.Span, node *commonpb.Node) *streampb.Span {
	return &streampb.Span{
		TraceID:       in.TraceId,
		SpanID:        in.SpanId,
		ParentSpanID:  in.ParentSpanId,
		ProcessName:   node.ServiceInfo.Name,
		OperationName: in.Name.Value,
		StartTime:     in.StartTime.Seconds,
		Duration:      int32((in.EndTime.Nanos - in.StartTime.Nanos) / 1000000),
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
