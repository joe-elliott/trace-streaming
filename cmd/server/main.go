package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	blergproto "github.com/joe-elliott/blerg/pkg/proto"
	"github.com/joe-elliott/blerg/pkg/streamer"
	"github.com/joe-elliott/blerg/pkg/util"
)

func main() {
	port := util.DefaultPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	server := grpc.NewServer()
	blergproto.RegisterSpanStreamServer(server, &streamer.Server{})

	server.Serve(lis)
}
