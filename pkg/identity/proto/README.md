protoc pkg/identity/proto/*.proto --go_out=plugins=grpc:.

protoc -I=. --proto_path=./third_party --go_out=plugins=grpc:. ./pkg/identity/proto/*.proto