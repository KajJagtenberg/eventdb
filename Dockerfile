FROM golang:1.16 AS build

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/eventflowdb cmd/eventflowdb/main.go

FROM debian:10

ENV DATA /data

COPY --from=build /src/bin/eventflowdb /usr/bin/eventflowdb

STOPSIGNAL SIGINT

ENV ENVIRONMENT=production

CMD ["eventflowdb"]