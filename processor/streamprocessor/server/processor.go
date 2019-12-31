package server

import "github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"

type StreamProcessor interface {
	AddSpanStreamer(s streamer.Streamer)
	AddTraceStreamer(s streamer.Streamer)
}

type StreamServer interface {
	Do() error
	Shutdown()
}
