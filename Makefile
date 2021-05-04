VERSION := $(shell cat constants/version)

db:
	mkdir -p data
	go run cmd/eventflowdb/main.go

ctl:
	go run cmd/eventflowctl/main.go

pb:
	protoc -I=${PWD} --go_out=${PWD}/ ${PWD}/proto/store.proto

clean:
	rm -rf data/*

test:
	go test ./...

build:
	docker build -t docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION) .

compile:
	go build cmd/eventflowdb/main.go

push: build
	docker tag docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION) docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:latest
	docker push docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:$(VERSION)
	docker push docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:latest