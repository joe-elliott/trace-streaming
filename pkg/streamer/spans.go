package streamer

import (
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type Spans struct {
	req    *blergpb.SpanRequest
	stream ClientStream
	spans  chan []*blergpb.Span
}

func NewSpans(req *blergpb.SpanRequest, stream ClientStream) *Spans {
	return &Spans{
		req:    req,
		stream: stream,
		spans:  make(chan []*blergpb.Span),
	}
}

func (s *Spans) Do() error {

	for spans := range s.spans {
		filtered := s.filterSpan(spans)

		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   filtered,
		})
	}

	return nil
}

func (s *Spans) ProcessBatch(spans []*blergpb.Span) {
	s.spans <- spans
}

func (s *Spans) Shutdown(spans []*blergpb.Span) {
	close(s.spans)
}

func (s *Spans) filterSpan(spans []*blergpb.Span) []*blergpb.Span {

	if len(s.req.ProcessName) > 0 || len(s.req.OperationName) > 0 {
		filtered := make([]*blergpb.Span, 0)

		for _, span := range spans {
			if len(s.req.ProcessName) > 0 && span.ProcessName == s.req.ProcessName {
				filtered = append(filtered, span)
				continue
			}

			if len(s.req.OperationName) > 0 && span.OperationName == s.req.OperationName {
				filtered = append(filtered, span)
				continue
			}
		}

		return filtered
	}

	return spans
}
