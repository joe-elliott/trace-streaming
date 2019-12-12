package streamprocessor

import (
	blergproto "github.com/joe-elliott/blerg/pkg/proto"
)

type tailer struct {
}

func (s *tailer) Tail(req *blergproto.StreamRequest, stream blergproto.SpanStream_TailServer) error {

	stream.Send(&blergproto.SpanResponse{
		Dropped: 32,
	})

	return nil
}
