package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streamer"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/traceql"
)

type WebsocketConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

type socketSender struct {
	ws *websocket.Conn
}

func (s *socketSender) Send(span *streampb.SpanResponse) error {
	return s.ws.WriteJSON(span)
}

type websocketServer struct {
	s      StreamProcessor
	cfg    WebsocketConfig
	server *http.Server
}

func NewWebsocket(s StreamProcessor, cfg WebsocketConfig) StreamServer {
	return &websocketServer{
		s:   s,
		cfg: cfg,
	}
}

func (w *websocketServer) Do() error {
	if !w.cfg.Enabled {
		return nil
	}

	http.HandleFunc("/v1/stream", w.stream)

	w.server = &http.Server{Addr: fmt.Sprintf(":%d", w.cfg.Port)}

	go func() {
		if err := w.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return nil
}

func (w *websocketServer) Shutdown() {
	if w.server != nil {
		w.server.Shutdown(context.TODO())
	}
}

func (s *websocketServer) stream(w http.ResponseWriter, r *http.Request) {
	socket := setupWebsocket(w, r)

	query := r.URL.Query()
	req := getStreamRequest(query)

	q, err := traceql.ParseExpr(req.Query)
	if err != nil {
		// todo: log or return error
		return
	}

	switch q.QueryType() {
	case traceql.QueryTypeBatchedSpans:
		tailer := streamer.NewBatchedSpans(q, int(req.RequestedRate), socket)
		s.s.AddTraceStreamer(tailer)
	case traceql.QueryTypeTraces:
		tailer := streamer.NewTraces(q, int(req.RequestedRate), socket)
		s.s.AddTraceStreamer(tailer)
	case traceql.QueryTypeMetrics:
		tailer := streamer.NewMetrics(q, int(req.RequestedRate), socket)
		s.s.AddSpanStreamer(tailer)
	case traceql.QueryTypeSpans:
		tailer := streamer.NewSpans(q, int(req.RequestedRate), socket)
		s.s.AddSpanStreamer(tailer)
	}
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
	query := getQueryParam(v, "q")

	return &streampb.StreamRequest{
		RequestedRate: int32(rate),
		Query:         query,
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
