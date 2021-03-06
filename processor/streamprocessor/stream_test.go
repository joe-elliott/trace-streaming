package streamprocessor

import (
	"context"
	"testing"

	tracepb "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exportertest"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"
)

func TestNewTraceProcessor(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	oCfg.GRPC.Enabled = false
	oCfg.Websocket.Enabled = false

	tp, err := NewTraceProcessor(nil, *oCfg)
	require.Error(t, oterr.ErrNilNextConsumer, err)
	require.Nil(t, tp)

	tp, err = NewTraceProcessor(exportertest.NewNopTraceExporter(), *oCfg)
	require.Nil(t, err)
	require.NotNil(t, tp)
}

// TestSpanProcessor_NilEmpty tests spans and attributes with nil/empty values
// do not cause any errors and no renaming occurs.
func TestSpanProcessor_NilEmpty(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	oCfg.GRPC.Enabled = false
	oCfg.Websocket.Enabled = false

	tp, err := factory.CreateTraceProcessor(zap.NewNop(), exportertest.NewNopTraceExporter(), oCfg)
	require.Nil(t, err)
	require.NotNil(t, tp)

	traceData := consumerdata.TraceData{
		Spans: []*tracepb.Span{
			nil,
			{
				Name:       &tracepb.TruncatableString{Value: "Nil Attributes"},
				Attributes: nil,
			},
			{
				Name:       &tracepb.TruncatableString{Value: "Empty Attributes"},
				Attributes: &tracepb.Span_Attributes{},
			},
			{
				Name: &tracepb.TruncatableString{Value: "Nil Attribute Map"},
				Attributes: &tracepb.Span_Attributes{
					AttributeMap: nil,
				},
			},
			{
				Name: &tracepb.TruncatableString{Value: "Empty Attribute Map"},
				Attributes: &tracepb.Span_Attributes{
					AttributeMap: map[string]*tracepb.AttributeValue{},
				},
			},
		},
	}
	assert.NoError(t, tp.ConsumeTraceData(context.Background(), traceData))
	assert.Equal(t, consumerdata.TraceData{
		Spans: []*tracepb.Span{
			nil,
			{
				Name:       &tracepb.TruncatableString{Value: "Nil Attributes"},
				Attributes: nil,
			},
			{
				Name:       &tracepb.TruncatableString{Value: "Empty Attributes"},
				Attributes: &tracepb.Span_Attributes{},
			},
			{
				Name: &tracepb.TruncatableString{Value: "Nil Attribute Map"},
				Attributes: &tracepb.Span_Attributes{
					AttributeMap: nil,
				},
			},
			{
				Name: &tracepb.TruncatableString{Value: "Empty Attribute Map"},
				Attributes: &tracepb.Span_Attributes{
					AttributeMap: map[string]*tracepb.AttributeValue{},
				},
			},
		},
	}, traceData)
}

// Common structure for the test cases.
type testCase struct {
	inputName        string
	inputAttributes  map[string]*tracepb.AttributeValue
	outputName       string
	outputAttributes map[string]*tracepb.AttributeValue
}

// runIndividualTestCase is the common logic of passing trace data through a configured attributes processor.
func runIndividualTestCase(t *testing.T, tt testCase, tp processor.TraceProcessor) {
	t.Run(tt.inputName, func(t *testing.T) {
		traceData := consumerdata.TraceData{
			Spans: []*tracepb.Span{
				{
					Name: &tracepb.TruncatableString{Value: tt.inputName},
					Attributes: &tracepb.Span_Attributes{
						AttributeMap: tt.inputAttributes,
					},
				},
			},
		}

		assert.NoError(t, tp.ConsumeTraceData(context.Background(), traceData))
		require.Equal(t, consumerdata.TraceData{
			Spans: []*tracepb.Span{
				{
					Name: &tracepb.TruncatableString{Value: tt.outputName},
					Attributes: &tracepb.Span_Attributes{
						AttributeMap: tt.outputAttributes,
					},
				},
			},
		}, traceData)
	})
}

// TestSpanProcessor_Values tests all possible value types.
func TestSpanProcessor_Values(t *testing.T) {
	testCases := []testCase{
		{
			inputName: "string type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
			},
			outputName: "string type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
			},
		},
		{
			inputName: "int type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
			},
			outputName: "int type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
			},
		},
		{
			inputName: "double type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
			},
			outputName: "double type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
			},
		},
		{
			inputName: "bool type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
			outputName: "bool type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
		},
		{
			inputName: "nil type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": nil,
			},
			outputName: "nil type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": nil,
			},
		},
		{
			inputName: "unknown type",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {},
			},
			outputName: "unknown type",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {},
			},
		},
	}

	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	oCfg.GRPC.Enabled = false
	oCfg.Websocket.Enabled = false

	tp, err := factory.CreateTraceProcessor(zap.NewNop(), exportertest.NewNopTraceExporter(), oCfg)
	require.Nil(t, err)
	require.NotNil(t, tp)
	for _, tc := range testCases {
		runIndividualTestCase(t, tc, tp)
	}
}

// TestSpanProcessor_MissingKeys tests that missing a key in an attribute map results in no span name changes.
func TestSpanProcessor_MissingKeys(t *testing.T) {
	testCases := []testCase{
		{
			inputName: "first keys missing",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
			outputName: "first keys missing",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
		},
		{
			inputName: "middle key missing",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
			outputName: "middle key missing",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
		},
		{
			inputName: "last key missing",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
			},
			outputName: "last key missing",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
			},
		},
		{
			inputName: "all keys exists",
			inputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
			outputName: "all keys exists",
			outputAttributes: map[string]*tracepb.AttributeValue{
				"key1": {
					Value: &tracepb.AttributeValue_StringValue{StringValue: &tracepb.TruncatableString{Value: "bob"}},
				},
				"key2": {
					Value: &tracepb.AttributeValue_IntValue{IntValue: 123},
				},
				"key3": {
					Value: &tracepb.AttributeValue_DoubleValue{DoubleValue: cast.ToFloat64(234.129312)},
				},
				"key4": {
					Value: &tracepb.AttributeValue_BoolValue{BoolValue: true},
				},
			},
		},
	}
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	oCfg.GRPC.Enabled = false
	oCfg.Websocket.Enabled = false

	tp, err := factory.CreateTraceProcessor(zap.NewNop(), exportertest.NewNopTraceExporter(), oCfg)
	require.Nil(t, err)
	require.NotNil(t, tp)
	for _, tc := range testCases {
		runIndividualTestCase(t, tc, tp)
	}
}
