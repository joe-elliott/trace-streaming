package main

import (
	"log"

	"github.com/open-telemetry/opentelemetry-collector/defaults"
	"github.com/open-telemetry/opentelemetry-collector/processor"
	"github.com/open-telemetry/opentelemetry-collector/service"

	"github.com/joe-elliott/blerg/pkg/streamprocessor"
)

func main() {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to run the service: %v", err)
		}
	}

	factories, err := defaults.Components()
	handleErr(err)

	// only need one processor for now.  can add more later
	factories.Processors, err = processor.Build(
		&streamprocessor.Factory{},
	)
	handleErr(err)

	svc := service.New(factories)
	err = svc.Start()
	handleErr(err)
}
