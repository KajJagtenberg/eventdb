build:
	docker build -t eventflowdb-ui:latest .

docker: build
	docker run --rm -it -p 3000:3000 eventflowdb-ui:latest