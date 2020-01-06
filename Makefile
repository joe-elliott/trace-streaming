.PHONY: tsp
tsp:
	GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/$(GOOS)/tsp $(BUILD_INFO) ./cmd/tsp

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