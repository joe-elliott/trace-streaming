package streamer

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"

type ClientStream interface {
	Send(s *streampb.SpanResponse) error
}

type Streamer interface {
	Shutdown()
	Do() error
	ProcessBatch([]*streampb.Span)
}
