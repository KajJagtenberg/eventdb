FROM golang:1.15-alpine AS build

WORKDIR /src

RUN apk update
RUN apk upgrade

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o eventflowdb .

FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add bash

WORKDIR /var/lib/eventdb

COPY --from=build /src/eventflowdb /bin/eventflowdb

CMD [ "eventflowdb" ]