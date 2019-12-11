package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	blergproto "github.com/joe-elliott/blerg/pkg/proto"
	"github.com/joe-elliott/blerg/pkg/util"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", util.DefaultPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	defer conn.Close()

	streamRequest := &blergproto.StreamRequest{}

	client := blergproto.NewSpanStreamClient(conn)
	stream, err := client.Tail(context.Background(), streamRequest)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("stream fail", err)
		}
		log.Println(resp)
	}

	log.Println("success!")
}
