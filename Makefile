proto:
	- protoc -I=. --proto_path=./third_party --go_out=plugins=grpc:. ./pkg/identity/proto/*.proto

lint:
	docker run --rm -v ${LOCAL_WORKSPACE_FOLDER}:/app -w /app golangci/golangci-lint:v1.41-alpine golangci-lint run ./... -v

infra:
	- docker-compose -f docker-compose-infra.yml up -d

infra-down:
	- docker-compose -f docker-compose-infra.yml down

test:
	go test -race -coverprofile=cover.out -covermode=atomic ./...
