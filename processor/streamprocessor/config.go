package streamprocessor

import (
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
)

// Config is the configuration for the stream processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
}
