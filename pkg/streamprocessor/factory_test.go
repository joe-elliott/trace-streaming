package streamprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector/config/configerror"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exportertest"
)

func TestFactory_Type(t *testing.T) {
	factory := &Factory{}
	assert.Equal(t, factory.Type(), typeStr)
}

func TestFactory_CreateDefaultConfig(t *testing.T) {
	factory := &Factory{}
	cfg := factory.CreateDefaultConfig()

	// Check the values of the default configuration.
	assert.NotNil(t, cfg)
	assert.Equal(t, typeStr, cfg.Type())
	assert.Equal(t, typeStr, cfg.Name())
}

func TestFactory_CreateTraceProcessor(t *testing.T) {
	factory := &Factory{}
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	// Name.FromAttributes field needs to be set for the configuration to be valid.
	oCfg.Rename.FromAttributes = []string{"test-key"}
	tp, err := factory.CreateTraceProcessor(zap.NewNop(), exportertest.NewNopTraceExporter(), oCfg)

	require.Nil(t, err)
	assert.NotNil(t, tp)
}

// TestFactory_CreateTraceProcessor_InvalidConfig ensures the default configuration
// returns an error.
func TestFactory_CreateTraceProcessor_InvalidConfig(t *testing.T) {
	factory := &Factory{}
	cfg := factory.CreateDefaultConfig()

	tp, err := factory.CreateTraceProcessor(zap.NewNop(), exportertest.NewNopTraceExporter(), cfg)
	require.Nil(t, tp)
	assert.Equal(t, err, errMissingRequiredField)
}

func TestFactory_CreateMetricProcessor(t *testing.T) {
	factory := &Factory{}
	cfg := factory.CreateDefaultConfig()

	mp, err := factory.CreateMetricsProcessor(zap.NewNop(), nil, cfg)
	require.Nil(t, mp)
	assert.Equal(t, err, configerror.ErrDataTypeIsNotSupported)
}
