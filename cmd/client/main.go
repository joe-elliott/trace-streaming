package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", 31234), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	defer conn.Close()

	traceReq := &streampb.TraceRequest{}

	client := streampb.NewSpanStreamClient(conn)
	stream, err := client.QueryTraces(context.Background(), traceReq)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("stream fail", err)
		}
		log.Println("----received thing----")
		log.Println(resp)
	}

	log.Println("success!")
}
