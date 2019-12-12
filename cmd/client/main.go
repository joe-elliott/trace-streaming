package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/joe-elliott/blerg/pkg/blergpb"
	"github.com/joe-elliott/blerg/pkg/util"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", util.DefaultPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	defer conn.Close()

	traceReq := &blergpb.TraceRequest{
		CrossesProcessBoundaries: false,
	}

	client := blergpb.NewSpanStreamClient(conn)
	stream, err := client.QueryTraces(context.Background(), traceReq)

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
