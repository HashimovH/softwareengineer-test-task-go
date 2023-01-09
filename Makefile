proto:
	protoc -I protos protos/tickets_service.proto --go_out=plugins=grpc:protos/tickets_service
