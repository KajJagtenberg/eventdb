run:
	go run .

build:
	go build .

clean:
	rm -f eventdb
	rm -rf *.bolt

docker:
	docker build -t eventdb .

compose_up:
	docker-compose up -d

compose_down:
	docker-compose down