# More exclusions can be added similar with: -not -path './testbed/*'
ALL_SRC := $(shell find . -name '*.go' \
                                -not -path './testbed/*' \
                                -type f | sort)

# All source code and documents. Used in spell check.
ALL_DOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                                -type f | sort)

# ALL_PKGS is used with 'go cover'
ALL_PKGS := $(shell go list $(sort $(dir $(ALL_SRC))))

GOTEST_OPT?= -race -timeout 30s
GOTEST_OPT_WITH_COVERAGE = $(GOTEST_OPT) -coverprofile=coverage.txt -covermode=atomic
GOTEST=go test
LINT=golangci-lint

.PHONY: tsp
tsp:
	GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/$(GOOS)/tsp $(BUILD_INFO) ./cmd/tsp

.PHONY: test
test:
	$(GOTEST) $(GOTEST_OPT) $(ALL_PKGS)

.PHONY: benchmark
benchmark:
	$(GOTEST) -bench=. -run=notests $(ALL_PKGS)

.PHONY: test-with-cover
test-with-cover:
	@echo Verifying that all packages have test files to count in coverage
	@scripts/check-test-files.sh $(subst github.com/open-telemetry/opentelemetry-collector/,./,$(ALL_PKGS))
	@echo pre-compiling tests
	@time go test -i $(ALL_PKGS)
	$(GOTEST) $(GOTEST_OPT_WITH_COVERAGE) $(ALL_PKGS)
	go tool cover -html=coverage.txt -o coverage.html

.PHONY: lint
lint:
	$(LINT) run

.PHONY: docker-component # Not intended to be used directly
docker-component: check-component
	GOOS=linux $(MAKE) $(COMPONENT)
	cp ./bin/linux/$(COMPONENT) ./cmd/$(COMPONENT)/
	docker build -t $(COMPONENT) ./cmd/$(COMPONENT)/
	rm ./cmd/$(COMPONENT)/$(COMPONENT)

.PHONY: docker-tsp
docker-tsp:
	COMPONENT=tsp $(MAKE) docker-component

.PHONY: check-component
check-component:
ifndef COMPONENT
	$(error COMPONENT variable was not defined)
endif