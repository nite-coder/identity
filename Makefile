proto:
	- protoc -I=. --proto_path=./third_party --go_out=plugins=grpc:. ./pkg/identity/proto/*.proto
