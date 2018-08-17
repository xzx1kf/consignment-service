build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/xzx1kf/consignment-service \
		proto/consignment/consignment.proto
	GOOS=linux GOARCH=arm go build
	docker build -t consignment-service .

run:
	docker run -d -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns consignment-service
