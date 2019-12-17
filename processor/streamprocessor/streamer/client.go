package streamer

import "github.com/joe-elliott/blerg/processor/streamprocessor/streampb"

type ClientStream interface {
	Send(s *streampb.SpanResponse) error
}
