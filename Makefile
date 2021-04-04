VERSION := $(shell cat constants/version)

eventflowdb:
	go run cmd/eventflowdb/main.go

eventflowctl:
	go run cmd/client/main.go

pb:
	protoc --proto_path=proto --go_out=plugins=grpc:api --go_opt=paths=source_relative proto/api.proto
	protoc --proto_path=proto --go_out=plugins=grpc:store --go_opt=paths=source_relative proto/store.proto

gql:
	go run github.com/99designs/gqlgen generate

build:
	docker build -t ghcr.io/kajjagtenberg/eventflowdb:$(VERSION) .

push: build
	docker tag ghcr.io/kajjagtenberg/eventflowdb:$(VERSION) ghcr.io/kajjagtenberg/eventflowdb:latest
	docker push ghcr.io/kajjagtenberg/eventflowdb:$(VERSION)
	docker push ghcr.io/kajjagtenberg/eventflowdb:latest

compose_up:
	docker-compose up -d --build

compose_down:
	docker-compose down