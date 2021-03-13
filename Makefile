run:
	go run cmd/server/main.go

tidy:
	go mod tidy

test:
	go test ./...

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