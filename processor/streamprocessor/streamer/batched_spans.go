package streamer

import (
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
	"go.uber.org/ratelimit"
)

type BatchedSpans struct {
	query   traceql.Query
	stream  ClientStream
	traces  chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewBatchedSpans(q traceql.Query, rate int, stream ClientStream) *BatchedSpans {

	return &BatchedSpans{
		query:   q,
		stream:  stream,
		traces:  make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *BatchedSpans) Do() error {
	for trace := range s.traces {
		trace = s.filterSpan(trace)

		if len(trace) == 0 {
			continue
		}

		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Type:    streampb.SpanResponse_SPANS,
			Spans:   trace,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *BatchedSpans) ProcessBatch(trace []*streampb.Span) {
	select {
	case s.traces <- trace:
	default:
		//todo: metric
	}
}

func (s *BatchedSpans) Shutdown() {
	close(s.traces)
}

func (s *BatchedSpans) filterSpan(trace []*streampb.Span) []*streampb.Span {
	filtered := make([]*streampb.Span, 0)

	for _, span := range trace {
		if s.query.MatchesSpanBatched(span, trace) {
			filtered = append(filtered, span)
		}
	}

	return filtered
}
