package streamer

import (
	"fmt"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
	"go.uber.org/ratelimit"
)

type Spans struct {
	query   traceql.Query
	stream  ClientStream
	spans   chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewSpans(q traceql.Query, rate int, stream ClientStream) *Spans {

	return &Spans{
		query:   q,
		stream:  stream,
		spans:   make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Spans) Do() error {
	for spans := range s.spans {
		spans = s.filterSpan(spans)

		if len(spans) == 0 {
			continue
		}

		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Spans:   spans,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Spans) ProcessBatch(spans []*streampb.Span) {
	select {
	case s.spans <- spans:
	default:
		fmt.Println("rate limited!")
	}
}

func (s *Spans) Shutdown() {
	close(s.spans)
}

func (s *Spans) filterSpan(spans []*streampb.Span) []*streampb.Span {
	filtered := make([]*streampb.Span, 0)

	for _, span := range spans {
		if s.query.MatchesSpan(span) {
			filtered = append(filtered, span)
		}
	}

	return filtered
}
