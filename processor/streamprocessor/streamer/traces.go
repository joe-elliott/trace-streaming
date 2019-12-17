package streamer

import (
	"fmt"

	"github.com/joe-elliott/blerg/processor/streamprocessor/streampb"
	"github.com/joe-elliott/blerg/processor/streamprocessor/util"
	"go.uber.org/ratelimit"
)

type Traces struct {
	req     *streampb.TraceRequest
	stream  ClientStream
	traces  chan []*streampb.Span
	limiter ratelimit.Limiter
}

func NewTraces(req *streampb.TraceRequest, stream ClientStream) *Traces {
	rate := util.DefaultRate
	if req.Params != nil && req.Params.RequestedRate != 0 {
		rate = int(req.Params.RequestedRate)
	}

	return &Traces{
		req:     req,
		stream:  stream,
		traces:  make(chan []*streampb.Span),
		limiter: ratelimit.New(rate),
	}
}

func (s *Traces) Do() error {

	for trace := range s.traces {
		s.stream.Send(&streampb.SpanResponse{
			Dropped: 0,
			Spans:   trace,
		})

		s.limiter.Take()
	}

	return nil
}

func (s *Traces) ProcessBatch(trace []*streampb.Span) {
	if !s.sendTrace(trace) {
		return
	}

	select {
	case s.traces <- trace:
	default:
		fmt.Println("rate limited!")
	}
}

func (s *Traces) Shutdown() {
	close(s.traces)
}

func (s *Traces) sendTrace(trace []*streampb.Span) bool {
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
