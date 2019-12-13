package streamer

import "github.com/joe-elliott/blerg/pkg/blergpb"

type ClientStream interface {
	Send(s *blergpb.SpanResponse) error
}
