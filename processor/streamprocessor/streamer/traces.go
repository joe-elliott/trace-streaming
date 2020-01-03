package streamer

import (
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
	"go.uber.org/ratelimit"
)

type Traces struct {
	query   traceql.Query
	stream  ClientStream
	traces  chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewTraces(q traceql.Query, rate int, stream ClientStream) *Traces {

	return &Traces{
		query:   q,
		stream:  stream,
		traces:  make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Traces) Do() error {
	for trace := range s.traces {
		trace = s.filterSpan(trace)

		if len(trace) == 0 {
			continue
		}

		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Type:    streampb.SpanResponse_TRACE,
			Spans:   trace,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Traces) ProcessBatch(trace []*streampb.Span) {
	select {
	case s.traces <- trace:
	default:
		//todo: metric
	}
}

func (s *Traces) Shutdown() {
	close(s.traces)
}

func (s *Traces) filterSpan(trace []*streampb.Span) []*streampb.Span {
	if s.query.MatchesTrace(trace) {
		return trace
	}

	return nil
}
