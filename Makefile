PLATFORM=linux
ARCHITECTURE=amd64
BINARY=server
FILES=./main

.PHONY: proto
proto:
	protoc -I protos protos/tickets_service.proto --go_out=plugins=grpc:protos/tickets_service

.PHONY: build
build:
	GOOS=$(PLATFORM) GOARCH=$(ARCHITECTURE) go build -ldflags="-s -w" -o build/$(BINARY) $(FILES)

.PHONY: build-docker-image
docker-build:
	@docker build . -t ticket-analysis-service-go

.PHONY: test
test: run-tests
	go test ./tests/...

.PHONY: lint
lint: run-lint
	golangci-lint run