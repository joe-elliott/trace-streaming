package streamprocessor

import (
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/server"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
)

// Config is the configuration for the stream processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
	GRPC                           server.GRPCConfig      `mapstructure:"grpc"`
	Websocket                      server.WebsocketConfig `mapstructure:"websocket"`
}
