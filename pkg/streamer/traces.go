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
		filtered := s.filterTraces(trace)

		s.stream.Send(&blergpb.SpanResponse{
			Dropped: 0,
			Spans:   filtered,
		})
	}

	return nil
}

func (s *Traces) ProcessBatch(trace []*blergpb.Span) {
	s.traces <- trace
}

func (s *Traces) Shutdown() {
	close(s.traces)
}

func (s *Traces) filterTraces(trace []*blergpb.Span) []*blergpb.Span {
	if !s.req.CrossesProcessBoundaries {
		var filtered []*blergpb.Span

		for _, span := range trace {
			if span.Parent != nil && span.Parent.ProcessName != span.ProcessName {
				filtered = append(filtered, span)
			}
		}

		return filtered
	}

	// unfiltered
	return trace
}
