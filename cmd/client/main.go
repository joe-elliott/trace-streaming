package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

var query string

func init() {
	flag.StringVar(&query, "query", "spans{}", "TraceQL query")
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf(":%d", 31234), grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to listen", err)
	}
	defer conn.Close()

	req := &streampb.StreamRequest{
		RequestedRate: 10,
		Query:         query,
	}

	client := streampb.NewSpanStreamClient(conn)
	stream, err := client.Query(context.Background(), req)

	if err != nil {
		log.Fatal("Failed to query", err)
	}

	log.Println("connected")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("stream fail", err)
		}
		fmt.Println("-----")
		b, _ := json.MarshalIndent(resp, "", "\t")
		fmt.Println(string(b))
	}

	log.Println("success!")
}
