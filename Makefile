VERSION := $(shell cat constants/version)

db:
	mkdir -p data
	go run cmd/eventflowdb/main.go

ctl:
	go run cmd/eventflowctl/main.go

pb:
	protoc --proto_path=proto --go_out=plugins=grpc:store --go_opt=paths=source_relative proto/store.proto

clean:
	rm -rf data/*

test:
	go test ./...

build:
	docker build -t docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION) .

push: build
	docker tag docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION) docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:latest
	docker push docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION)
	docker push docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:latest

compose_up:
	docker-compose up -d --build

compose_down:
	docker-compose down