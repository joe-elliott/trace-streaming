package server

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"

type StreamProcessor interface {
	AddSpanStreamer(s *streamer.Spans)
	AddTraceStreamer(s *streamer.Traces)
}
