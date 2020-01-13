package main

import (
	"log"

	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/open-telemetry/opentelemetry-collector/exporter"
	"github.com/open-telemetry/opentelemetry-collector/processor"
	"github.com/open-telemetry/opentelemetry-collector/service"

	"github.com/joe-elliott/trace-streaming/exporter/nativeexporter"
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor"
)

func main() {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to run the service: %v", err)
		}
	}

	factories, err := defaults.Components()
	handleErr(err)

	info := service.ApplicationStartInfo{
		ExeName:  "tsp",
		LongName: "trace stream processing",
		Version:  "x.x.x.x",
		GitHash:  "",
	}

	// processors
	customProcessors, err := processor.Build(
		&streamprocessor.Factory{},
	)
	handleErr(err)

	for k, v := range customProcessors {
		factories.Processors[k] = v
	}

	// exporters
	customExporters, err := exporter.Build(
		&nativeexporter.Factory{},
	)
	handleErr(err)

	for k, v := range customExporters {
		factories.Exporters[k] = v
	}

	// start
	svc, err := service.New(factories, info)
	handleErr(err)

	err = svc.Start()
	handleErr(err)
}
