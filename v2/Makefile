server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

pb:
	protoc --proto_path=proto --go_out=plugins=grpc:cluster --go_opt=paths=source_relative proto/fsm.proto
	protoc --proto_path=proto --go_out=plugins=grpc:api --go_opt=paths=source_relative proto/api.proto
	protoc --proto_path=proto --go_out=plugins=grpc:persistence --go_opt=paths=source_relative proto/persistence.proto

gql:
	go run github.com/99designs/gqlgen generate

build:
	docker build -t eventflowdb:latest .

compose_up:
	docker-compose up -d --build

compose_down:
	docker-compose down