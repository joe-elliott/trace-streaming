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
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter/backends"
	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter/backends/local"
	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter/batch"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/open-telemetry/opentelemetry-collector/consumer/consumerdata"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"github.com/open-telemetry/opentelemetry-collector/exporter/exporterhelper"
)

type nativeExporter struct {
	batcher batch.Batcher
	backend backends.Writer

	batches []batch.Batch

	logger *zap.Logger
	cfg    *Config
}

// NewTraceExporter creates an exporter.TraceExporter that just drops the
// received data and logs debugging messages.
func NewTraceExporter(config configmodels.Exporter, logger *zap.Logger) (exporter.TraceExporter, error) {
	dir, err := ioutil.TempDir("", "tsp-test")
	if err != nil {
		return nil, err
	}

	e := &nativeExporter{
		batcher: batch.NewBatcher(),
		backend: local.NewBackend(dir),
		logger:  logger,
		cfg:     config.(*Config),
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
	ticker := time.NewTicker(e.cfg.BlockDuration)

	for {
		batch, err := e.batcher.Cut()

		if err != nil {
			e.logger.Error("Error cutting batch.", zap.Error(err))
		}

		e.batches = append(e.batches, batch)

		<-ticker.C
	}
}

func (e *nativeExporter) flushBatches() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		<-ticker.C

		// todo : shared structure needs locking?  this seems generally like a poor way to do this.  find a better way
		if len(e.batches) == 0 {
			continue
		}

		b := e.batches[0]

		idxBytes, err := b.Index()
		if err != nil {
			e.logger.Error("Error getting index.", zap.Error(err))
			continue
		}

		traceBytes, err := b.Traces()
		if err != nil {
			e.logger.Error("Error getting traces.", zap.Error(err))
			continue
		}

		filterBytes, err := b.BloomFilter()
		if err != nil {
			e.logger.Error("Error getting filter.", zap.Error(err))
			continue
		}

		err = e.backend.Write(idxBytes, traceBytes, filterBytes)
		if err != nil {
			e.logger.Error("Error flushing.", zap.Error(err))
			continue
		}
	}
}
