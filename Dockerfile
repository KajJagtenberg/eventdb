FROM golang:1.16 AS build

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/eventflowdb cmd/eventflowdb/main.go

FROM debian

ENV DATA /data
ENV TLS_CERT_FILE /certs/cert.pem
ENV TLS_KEY_FILE /certs/key.pem

RUN mkdir $DATA/{fsm,snapshots} -p

COPY --from=build /src/bin/eventflowdb /usr/bin/eventflowdb

STOPSIGNAL SIGINT

ENV ENVIRONMENT=production

CMD ["eventflowdb"]