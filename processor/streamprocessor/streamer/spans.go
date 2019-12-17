package streamer

import (
	"fmt"

	"github.com/joe-elliott/blerg/processor/streamprocessor/streampb"
	"github.com/joe-elliott/blerg/processor/streamprocessor/util"
	"go.uber.org/ratelimit"
)

type Spans struct {
	req     *streampb.SpanRequest
	stream  ClientStream
	spans   chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewSpans(req *streampb.SpanRequest, stream ClientStream) *Spans {
	rate := util.DefaultRate
	if req.Params.RequestedRate != 0 {
		rate = int(req.Params.RequestedRate)
	}

	return &Spans{
		req:     req,
		stream:  stream,
		spans:   make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Spans) Do() error {

	for spans := range s.spans {
		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Spans:   spans,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Spans) ProcessBatch(spans []*streampb.Span) {
	filtered := s.filterSpan(spans)

	if len(filtered) == 0 {
		return
	}

	select {
	case s.spans <- filtered:
	default:
		fmt.Println("rate limited!")
	}
}

func (s *Spans) Shutdown() {
	close(s.spans)
}

func (s *Spans) filterSpan(spans []*streampb.Span) []*streampb.Span {

	if len(s.req.ProcessName) > 0 || len(s.req.OperationName) > 0 || s.req.MinDuration > 0 {
		filtered := make([]*streampb.Span, 0)

		for _, span := range spans {
			if (len(s.req.ProcessName) == 0 || span.ProcessName == s.req.ProcessName) &&
				(len(s.req.OperationName) == 0 || span.OperationName == s.req.OperationName) &&
				(s.req.MinDuration == 0 || span.Duration >= s.req.MinDuration) {
				filtered = append(filtered, span)
				continue
			}
		}

		return filtered
	}

	return spans
}
