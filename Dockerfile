FROM golang:1.16-alpine AS build

RUN apk add g++

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/server cmd/server/main.go

FROM alpine:3.8

WORKDIR /var/lib/eventflowdb

RUN mkdir data -p

COPY --from=build /src/bin/server /bin/eventflowdb

CMD ["eventflowdb"]