package streamer

import (
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
	"go.uber.org/ratelimit"
)

type Metrics struct {
	query   traceql.Query
	stream  ClientStream
	traces  chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewMetrics(q traceql.Query, rate int, stream ClientStream) *Metrics {

	return &Metrics{
		query:   q,
		stream:  stream,
		traces:  make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Metrics) Do() error {
	for trace := range s.traces {
		trace = s.filterSpan(trace)

		if len(trace) == 0 {
			continue
		}

		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Spans:   trace,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Metrics) ProcessBatch(trace []*streampb.Span) {
	select {
	case s.traces <- trace:
	default:
		//todo: metric
	}
}

func (s *Metrics) Shutdown() {
	close(s.traces)
}

func (s *Metrics) filterSpan(spans []*streampb.Span) []*streampb.Span {
	filtered := make([]*streampb.Span, 0)

	for _, span := range spans {
		if s.query.MatchesSpan(span) {
			// jpe calc metrics
		}
	}

	return filtered
}
