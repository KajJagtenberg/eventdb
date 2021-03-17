run:
	air

tidy:
	go mod tidy

test:
	go test ./... -v

clean:
	rm -f eventflowdb
	rm -rf *.bolt

generate:
	go generate ./...

build:
	docker build -t eventflowdb .

compose_up:
	docker-compose up -d

compose_down:
	docker-compose down

client:
	go run cmd/client/main.go

protobuf:
	protoc -I=protos --go_out=./gen protos/*.proto