FROM golang:1.16-alpine AS build

RUN apk add g++

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test ./...

RUN go build -o bin/eventflowdb cmd/eventflowdb/main.go

FROM alpine:3.8

ENV DATA /data
ENV DATA /certs/cert.pem
ENV DATA /certs/key.pem

RUN mkdir $DATA

COPY --from=build /src/bin/eventflowdb /usr/bin/eventflowdb

STOPSIGNAL SIGINT

CMD ["eventflowdb"]