package streamer

import "github.com/joe-elliott/blerg/processor/streamprocessor/blergpb"

type ClientStream interface {
	Send(s *blergpb.SpanResponse) error
}
