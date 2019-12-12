package streamprocessor

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector/consumer"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/oterr"
	"github.com/open-telemetry/opentelemetry-collector/processor"
)

type streamProcessor struct {
	nextConsumer consumer.TraceConsumer
	config       Config
}

// NewTraceProcessor returns the span processor.
func NewTraceProcessor(nextConsumer consumer.TraceConsumer, config Config) (processor.TraceProcessor, error) {
	if nextConsumer == nil {
		return nil, oterr.ErrNilNextConsumer
	}

	sp := &streamProcessor{
		nextConsumer: nextConsumer,
		config:       config,
	}

	return sp, nil
}

func (sp *streamProcessor) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {
	return sp.nextConsumer.ConsumeTraceData(ctx, td)
}

func (sp *streamProcessor) GetCapabilities() processor.Capabilities {
	return processor.Capabilities{MutatesConsumedData: false}
}
