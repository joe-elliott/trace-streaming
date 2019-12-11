package streamer

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	blergproto "github.com/joe-elliott/blerg/pkg/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestStream(t *testing.T) {
	port := getAvailablePort(t)
	_, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	require.NoError(t, err)

	server := grpc.NewServer()
	blergproto.RegisterSpanStreamServer(server, &Server{})
}

func getAvailableLocalAddress(t *testing.T) string {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to get a free local port: %v", err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

func getAvailablePort(t *testing.T) uint16 {
	endpoint := getAvailableLocalAddress(t)
	_, port, err := net.SplitHostPort(endpoint)
	require.NoError(t, err)

	portInt, err := strconv.Atoi(port)
	require.NoError(t, err)

	return uint16(portInt)
}
