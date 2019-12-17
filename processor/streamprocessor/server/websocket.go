package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/util"
)

type socketSender struct {
	ws *websocket.Conn
}

func (s *socketSender) Send(span *streampb.SpanResponse) error {
	return s.ws.WriteJSON(span)
}

type websocketServer struct {
	s StreamProcessor
}

func DoWebsocket(s StreamProcessor) error {
	w := websocketServer{
		s: s,
	}

	http.HandleFunc("/v1/stream/traces", w.streamTraces)
	http.HandleFunc("/v1/stream/spans", w.streamSpans)

	go http.ListenAndServe(fmt.Sprintf(":%d", util.DefaultHTTPPort), nil)

	return nil
}

func (s *websocketServer) streamTraces(w http.ResponseWriter, r *http.Request) {
	socket := setupWebsocket(w, r)

	query := r.URL.Query()

	tailer := streamer.NewTraces(&streampb.TraceRequest{
		Params:      getStreamRequest(query),
		ProcessName: getQueryParam(query, "processName"),
	}, socket)

	s.s.AddTraceStreamer(tailer)
}

func (s *websocketServer) streamSpans(w http.ResponseWriter, r *http.Request) {
	socket := setupWebsocket(w, r)

	query := r.URL.Query()

	tailer := streamer.NewSpans(&streampb.SpanRequest{
		Params:        getStreamRequest(query),
		ProcessName:   getQueryParam(query, "processName"),
		OperationName: getQueryParam(query, "operationName"),
		MinDuration:   int32(getQueryParamInt(query, "minDuration")),
	}, socket)

	s.s.AddSpanStreamer(tailer)
}

// utility
func getQueryParam(v url.Values, name string) string {
	value, ok := v[name]

	if ok && len(value) > 0 {
		return value[0]
	}

	return ""
}

func getQueryParamInt(v url.Values, name string) int {
	value, ok := v[name]

	ret := 0

	if ok && len(value) > 0 {
		ret, _ = strconv.Atoi(value[0])
	}

	return ret
}

func getStreamRequest(v url.Values) *streampb.StreamRequest {
	rate := getQueryParamInt(v, "rate")

	return &streampb.StreamRequest{
		RequestedRate: int32(rate),
	}
}

func setupWebsocket(w http.ResponseWriter, r *http.Request) *socketSender {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// helpful log statement to show connections
	log.Println("Client Connected")

	return &socketSender{
		ws: ws,
	}
}
