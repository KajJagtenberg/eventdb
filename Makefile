VERSION := $(shell cat constants/version)

db:
	mkdir -p data
	go run cmd/eventflowdb/main.go

compile:
	go build cmd/eventflowdb/main.go

pb:
	protoc -I=${PWD} --go_out=. ${PWD}/proto/store.proto
	protoc -I=${PWD} --go_out=.  --go-grpc_out=. ${PWD}/proto/api.proto


tidy:
	go mod tidy

clean:
	rm -rf data/*

test:
	go test ./...

build:
	docker build -t kajjagtenberg/eventflowdb:$(VERSION) .

push: build
	docker push kajjagtenberg/eventflowdb:$(VERSION)