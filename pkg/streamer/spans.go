package streamer

import (
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type Spans struct {
	req    *blergpb.SpanRequest
	stream blergpb.SpanStream_QuerySpansServer
	spans  chan []*blergpb.Span
}

func NewSpans(req *blergpb.SpanRequest, stream blergpb.SpanStream_QuerySpansServer) *Spans {
	return &Spans{
		req:    req,
		stream: stream,
		spans:  make(chan []*blergpb.Span),
	}
}

func (s *Spans) Do() error {

	for spans := range s.spans {
		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   spans,
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
