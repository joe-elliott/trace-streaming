package streamer

import (
	"time"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
)

type Metrics struct {
	query  traceql.Query
	stream ClientStream
	traces chan []*streampb.Span
	rate   time.Duration
}

func NewMetrics(q traceql.Query, rate int, stream ClientStream) *Metrics {

	return &Metrics{
		query:  q,
		stream: stream,
		traces: make(chan []*streampb.Span),
		rate:   time.Duration(rate) * time.Second,
	}
}

func (s *Metrics) Do() error {
	go func() {
		// todo: shutdown cleanly
		for {
			// todo: totally a race condition here
			metrics := s.query.Aggregate(nil, true)

			s.stream.Send(&streampb.SpanResponse{
				Dropped: 0,
				Type:    streampb.SpanResponse_METRICS,
				Metrics: &streampb.Metric{
					V: metrics,
					T: time.Now().Unix(),
				},
			})

			time.Sleep(s.rate)
		}
	}()

	for trace := range s.traces {
		for _, span := range trace {
			s.query.Aggregate(span, false)
		}
	}

	return nil
}

func (s *Metrics) ProcessBatch(trace []*streampb.Span) {
	select {
	case s.traces <- trace:
	default:
		//todo: metric
	}
}

func (s *Metrics) Shutdown() {
	close(s.traces)
}
