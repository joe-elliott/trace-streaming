package streamprocessor

import (
	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type spanTailer struct {
	spans chan []*tracepb.Span
}

func newSpanTailer() *spanTailer {
	return &spanTailer{
		spans: make(chan []*tracepb.Span),
	}
}

func (s *spanTailer) Tail(req *blergpb.StreamRequest, stream blergpb.SpanStream_TailServer) error {

	for spans := range s.spans {
		blergSpans := make([]*blergpb.Span, len(spans))

		for i, span := range spans {
			blergSpan := spanToSpan(span)
			blergSpans[i] = blergSpan
		}

		stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   blergSpans,
		})
	}

	return nil
}

func (s *spanTailer) processBatch(spans []*tracepb.Span) {
	s.spans <- spans
}

func (s *spanTailer) shutdown(spans []*tracepb.Span) {
	close(s.spans)
}

func spanToSpan(in *tracepb.Span) *blergpb.Span {

	return &blergpb.Span{
		OperationName: in.Name.Value,
		StartTime:     in.StartTime.Seconds,
		Duration:      in.EndTime.Seconds - in.StartTime.Seconds,
	}
}
