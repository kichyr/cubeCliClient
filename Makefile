.PHONY: docker-ci
docker-ci:
	@docker build --no-cache --network="host" -t cubeclient-ci:0.0.0 -f docker/ci.dockerfile .

.PHONY: test-local
test-local:
	@go test -count=1 ./pkg/...
	@go test -count=1 ./test/...
	@go test -count=1 ./cmd/...
	make functional-tests

functional-tests: build build-test-server
	pytest -s test/functional_test.py
	make clear

.PHONY: lint
lint:
	@golangci-lint run --config .golangci.yaml

.PHONY: start-test-server
start-test-server:
	@go run ./test/testserver/main.go

.PHONY: build
build:
	@go build github.com/kichyr/cubeCliClient/cmd/cubeclient

.PHONY: build-test-server
build-test-server:
	@go build github.com/kichyr/cubeCliClient/cmd/testserver

.PHONY: clear
clear:
	rm cubeclient testserver