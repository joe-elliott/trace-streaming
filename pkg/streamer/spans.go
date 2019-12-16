package streamer

import (
	"fmt"

	"github.com/joe-elliott/blerg/pkg/blergpb"
	"github.com/joe-elliott/blerg/pkg/util"
	"go.uber.org/ratelimit"
)

type Spans struct {
	req     *blergpb.SpanRequest
	stream  ClientStream
	spans   chan []*blergpb.Span
	limiter ratelimit.Limiter
}

func NewSpans(req *blergpb.SpanRequest, stream ClientStream) *Spans {
	rate := util.DefaultRate
	if req.Params.RequestedRate != 0 {
		rate = int(req.Params.RequestedRate)
	}

	return &Spans{
		req:     req,
		stream:  stream,
		spans:   make(chan []*blergpb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Spans) Do() error {

	for spans := range s.spans {
		filtered := s.filterSpan(spans)

		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   filtered,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Spans) ProcessBatch(spans []*blergpb.Span) {
	select {
	case s.spans <- spans:
	default:
		fmt.Println("rate limited!")
	}
}

func (s *Spans) Shutdown(spans []*blergpb.Span) {
	close(s.spans)
}

func (s *Spans) filterSpan(spans []*blergpb.Span) []*blergpb.Span {

	if len(s.req.ProcessName) > 0 || len(s.req.OperationName) > 0 {
		filtered := make([]*blergpb.Span, 0)

		for _, span := range spans {
			if (len(s.req.ProcessName) == 0 || span.ProcessName == s.req.ProcessName) &&
				(len(s.req.OperationName) == 0 || span.OperationName == s.req.OperationName) &&
				(s.req.MinDuration > 0 || span.Duration >= int64(s.req.MinDuration)) {
				filtered = append(filtered, span)
				continue
			}
		}

		return filtered
	}

	return spans
}
