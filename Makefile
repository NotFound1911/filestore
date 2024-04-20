# 根据 .proto 生成 .go   .mock
generate:
	mkdir -p ./api/proto/gen
	# account
	protoc --go_out=./api/proto/gen --go-grpc_out=./api/proto/gen ./api/proto/account/v1/account.proto
	# file_manager
	protoc --go_out=./api/proto/gen --go-grpc_out=./api/proto/gen ./api/proto/file_manager/v1/file_manager.proto
mock:
	# account
	mockgen -source=./api/proto/gen/account/v1/account_grpc.pb.go -package=accountmocks -destination=./api/proto/gen/account/v1/mocks/account_grpc.mock.go
	# file_manager
	mockgen -source=./api/proto/gen/file_manager/v1/file_manager_grpc.pb.go -package=filemanagermocks -destination=./api/proto/gen/filemanager/v1/mocks/filemanager_grpc.mock.go
clean:
	rm -rf ./api/proto/gen/*

start:
	# todo
	go run ./app/account/main.go
	go run ./app/apigw/main.go
	go run ./app/upload/main.go

# Default target (run 'make' without arguments)
all: clean generate mock