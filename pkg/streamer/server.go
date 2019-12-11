package streamer

import (
	blergproto "github.com/joe-elliott/blerg/pkg/proto"
)

type Server struct {
}

func (s *Server) Tail(req *blergproto.StreamRequest, stream blergproto.SpanStream_TailServer) error {
	return nil
}
