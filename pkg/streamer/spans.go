package streamer

import (
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type Spans struct {
	req    *blergpb.StreamRequest
	stream blergpb.SpanStream_TailServer
	spans  chan []*tracepb.Span
}

func NewSpans(req *blergpb.StreamRequest, stream blergpb.SpanStream_TailServer) *Spans {
	return &Spans{
		req:    req,
		stream: stream,
		spans:  make(chan []*tracepb.Span),
	}
}

func (s *Spans) Do() error {

	for spans := range s.spans {
		blergSpans := make([]*blergpb.Span, len(spans))

		for i, span := range spans {
			blergSpan := spanToSpan(span)
			blergSpans[i] = blergSpan
		}

		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   blergSpans,
		})
	}

	return nil
}

func (s *Spans) ProcessBatch(spans []*tracepb.Span) {
	s.spans <- spans
}

func (s *Spans) Shutdown(spans []*tracepb.Span) {
	close(s.spans)
}
