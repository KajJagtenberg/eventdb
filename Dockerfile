FROM golang:1.15-alpine AS build

WORKDIR /src

RUN apk update
RUN apk upgrade

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go install github.com/br0xen/boltbrowser

COPY . .

RUN go install .

FROM node:alpine AS webui

WORKDIR /src

COPY webui/package.json .
COPY webui/yarn.lock .

RUN yarn

COPY webui/ .

RUN yarn export

FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add bash

WORKDIR /var/lib/eventdb

COPY --from=build /go/bin/* /bin/
COPY --from=webui /src/out webui

CMD [ "eventdb" ]