run:
	go run .

build:
	go build -o eventflowdb main.go

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -f eventflowdb
	rm -rf *.bolt

docker:
	docker build -t eventflowdb .

compose_up:
	docker-compose up -d

compose_down:
	docker-compose down