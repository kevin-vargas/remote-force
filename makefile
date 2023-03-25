gen:
	protoc --proto_path=proto ./proto/*.proto --go_out=pb --go-grpc_out=pb
build:
# go build -o ./bin/client ./client
	@echo "Build server"
	@go build -o ./bin/server ./server