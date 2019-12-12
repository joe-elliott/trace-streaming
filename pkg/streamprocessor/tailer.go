package streamprocessor

import (
	"github.com/joe-elliott/blerg/pkg/blergpb"
)

type tailer struct {
}

func (s *tailer) Tail(req *blergpb.StreamRequest, stream blergpb.SpanStream_TailServer) error {

	stream.Send(&blergpb.SpanResponse{
		Dropped: 32,
	})

	return nil
}
