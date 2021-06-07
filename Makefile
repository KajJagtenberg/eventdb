VERSION := $(shell cat constants/version)

db:
	mkdir -p data
	go run cmd/eventflowdb/main.go

compile:
	go build cmd/eventflowdb/main.go

ctl:
	go run cmd/eventflowctl/main.go

pb:
	#  protoc -I=${PWD} --go_out=${PWD} ${PWD}/proto/store.proto
	protoc -I=${PWD} --go_out=. ${PWD}/proto/store.proto
	protoc -I=${PWD} --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto

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