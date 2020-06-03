.PHONY: docker-ci
docker-ci:

	@docker build --no-cache -t cubeclient-ci:0.0.0 -f docker/ci.dockerfile .

.PHONY: test
test: docker-ci
	@docker run cubeclient-ci:0.0.0

.PHONY: test-local
test-local:
	@go test -count=1 ./pkg/...
	@go test -count=1 ./test/...
	@go test -count=1 ./cmd/...
	make functional-tests


.PHONY: functional-tests
functional-tests: build build-test-server
	pytest -s test/functional_test.py
	make clear

.PHONY: lint
lint:
	@golangci-lint run --config .golangci.yaml

.PHONY: start-test-server
start-test-server:
	@go run ./cmd/testserver/main.go 8091

.PHONY: build
build:
	@go build github.com/kichyr/cubeCliClient/cmd/cubeclient

.PHONY: build-test-server
build-test-server:
	@go build github.com/kichyr/cubeCliClient/cmd/testserver

.PHONY: clear
clear:
	rm cubeclient testserver