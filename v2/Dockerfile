FROM node:14-alpine AS webui

WORKDIR /src

COPY webui/package.json .
COPY webui/yarn.lock .

RUN yarn

COPY webui/ .

RUN yarn export

FROM golang:1.16-alpine AS build

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux go build -o bin/server -ldflags="-s -w" cmd/server/main.go

FROM alpine:3.8

COPY --from=build /src/bin/server /bin/eventflowdb

COPY --from=webui /src/out/ webui/out

CMD ["eventflowdb"]