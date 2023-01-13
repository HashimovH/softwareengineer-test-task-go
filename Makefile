proto:
	protoc -I protos protos/tickets_service.proto --go_out=plugins=grpc:protos/tickets_service


.PHONY: build-docker-image
build:
    @docker build . -t ticket-analysis-service-go

.PHONY: run-docker
run-docker: build-docker-image
    @docker run -p 8000:8000 ticket-analysis-service-go

.PHONY: test
test: run-tests
    go test ./tests/...