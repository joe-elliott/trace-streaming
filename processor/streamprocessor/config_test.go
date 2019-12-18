package streamprocessor

import (
	"path"
	"testing"

	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/server"
	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
)

func TestLoadConfig(t *testing.T) {
	factories, err := config.ExampleComponents()
	assert.NoError(t, err)

	factory := &Factory{}
	factories.Processors[typeStr] = factory

	config, err := config.LoadConfigFile(t, path.Join(".", "testdata", "config.yaml"), factories)

	assert.Nil(t, err)
	assert.NotNil(t, config)

	p0 := config.Processors["stream"]
	assert.Equal(t, p0, &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: typeStr,
			NameVal: "stream",
		},
		GRPC: server.GRPCConfig{
			Port:    1111,
			Enabled: true,
		},
		Websocket: server.WebsocketConfig{
			Port:    1234,
			Enabled: false,
		},
	})

	p1 := config.Processors["stream/customname"]
	assert.Equal(t, p1, &Config{
		ProcessorSettings: configmodels.ProcessorSettings{
			TypeVal: typeStr,
			NameVal: "stream/customname",
		},
		GRPC: server.GRPCConfig{
			Port:    31234,
			Enabled: true,
		},
		Websocket: server.WebsocketConfig{
			Port:    31235,
			Enabled: true,
		},
	})
}
