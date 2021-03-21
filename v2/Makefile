server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

pb:
	protoc --proto_path=proto --go_out=plugins=grpc:store --go_opt=paths=source_relative proto/streams.proto
	protoc --proto_path=proto --go_out=plugins=grpc:cluster --go_opt=paths=source_relative proto/cluster.proto

gql:
	go run github.com/99designs/gqlgen generate

build:
	docker build -t eventflowdb:latest .

compose_up:
	docker-compose up -d --build

compose_down:
	docker-compose down