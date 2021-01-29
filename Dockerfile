FROM golang:1.15-alpine AS build

WORKDIR /src

RUN apk update
RUN apk upgrade

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go install github.com/br0xen/boltbrowser

COPY . .

RUN go build -o eventdb .

FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add bash

WORKDIR /var/lib/eventdb

COPY --from=build /src/eventdb /bin/eventdb
COPY --from=build /go/bin/* /bin/

CMD [ "eventdb" ]