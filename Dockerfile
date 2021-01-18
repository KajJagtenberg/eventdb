FROM golang:1.15-alpine AS build

WORKDIR /src

RUN apk update
RUN apk upgrade
RUN apk add gcc g++

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build .

FROM alpine

RUN apk update
RUN apk upgrade

WORKDIR /var/lib/eventdb

COPY --from=build /src/eventdb /bin/eventdb

CMD [ "eventdb" ]