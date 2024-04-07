# 根据 .proto 生成 .go   .mock
generate:
	mkdir -p ./api/proto/gen
	# account
	protoc --go_out=./api/proto/gen --go-grpc_out=./api/proto/gen ./api/proto/account/v1/account.proto

mock:
	# account
	mockgen -source=./api/proto/gen/account/v1/account_grpc.pb.go -package=accountmocks -destination=./api/proto/gen/account/v1/mocks/account_grpc.mock.go

clean:
	rm -rf ./api/proto/gen/*
# Default target (run 'make' without arguments)
all: clean generate mock