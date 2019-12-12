package streamer

import (
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type Traces struct {
	req    *blergpb.TraceRequest
	stream blergpb.SpanStream_QueryTracesServer
	traces chan []*blergpb.Span
}

func NewTraces(req *blergpb.TraceRequest, stream blergpb.SpanStream_QueryTracesServer) *Traces {
	return &Traces{
		req:    req,
		stream: stream,
		traces: make(chan []*blergpb.Span),
	}
}

func (s *Traces) Do() error {

	for trace := range s.traces {
		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   trace,
		})
	}

	return nil
}

func (s *Traces) ProcessBatch(trace []*blergpb.Span) {
	s.traces <- trace
}

func (s *Traces) Shutdown(spans []*blergpb.Span) {
	close(s.traces)
}
