FROM golang:1.15-alpine AS build

WORKDIR /src

RUN apk add g++

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test ./...
RUN go build -o eventflowdb .

FROM alpine

RUN apk add bash

WORKDIR /var/lib/eventflowdb

COPY --from=build /src/eventflowdb /bin/eventflowdb

CMD [ "eventflowdb" ]