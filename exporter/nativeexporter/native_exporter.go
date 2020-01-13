// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nativeexporter

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter/batch"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exporterhelper"
)

type nativeExporter struct {
	batcher batch.Batcher

	logger *zap.Logger
}

// NewTraceExporter creates an exporter.TraceExporter that just drops the
// received data and logs debugging messages.
func NewTraceExporter(config configmodels.Exporter, logger *zap.Logger) (exporter.TraceExporter, error) {
	e := &nativeExporter{
		batcher: batch.NewBatcher(),
		logger:  logger,
	}

	go e.cutTraces()

	return exporterhelper.NewTraceExporter(
		config,
		e.consumeTrace,
		exporterhelper.WithTracing(true),
		exporterhelper.WithMetrics(true),
		exporterhelper.WithShutdown(logger.Sync), // todo: flush on shutdown
	)
}

func (e *nativeExporter) consumeTrace(ctx context.Context, td consumerdata.TraceData) (int, error) {
	e.batcher.Store(td)

	return 0, nil
}

func (e *nativeExporter) cutTraces() {
	ticker := time.NewTicker(1 * time.Hour)

	for {
		_, err := e.batcher.Cut() // todo: write this somewhere

		if err != nil {
			e.logger.Error("Error cutting batch.", zap.Error(err))
		}

		<-ticker.C
	}

}
