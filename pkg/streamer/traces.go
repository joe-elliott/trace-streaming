package streamer

import (
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type Traces struct {
	req    *blergpb.TraceRequest
	stream ClientStream
	traces chan []*blergpb.Span
}

func NewTraces(req *blergpb.TraceRequest, stream ClientStream) *Traces {
	return &Traces{
		req:    req,
		stream: stream,
		traces: make(chan []*blergpb.Span),
	}
}

func (s *Traces) Do() error {

	for trace := range s.traces {
		if !s.sendTrace(trace) {
			continue
		}

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

func (s *Traces) Shutdown() {
	close(s.traces)
}

func (s *Traces) sendTrace(trace []*blergpb.Span) bool {
	if len(s.req.ProcessName) > 0 {
		for _, span := range trace {
			if span.ProcessName == s.req.ProcessName {
				return true
			}
		}

		return false
	}

	// unfiltered
	return true
}
