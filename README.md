# EventflowDB

EventflowDB is a database designed with Event Sourcing in mind.

- [EventflowDB](#eventflowdb)
    - [Features](#features)
    - [Prerequisites](#prerequisites)
    - [Installing](#installing)
  - [Configuration](#configuration)
  - [Usage](#usage)
  - [Example](#example)
  - [Versioning](#versioning)
  - [Roadmap](#roadmap)
  - [Contributions](#contributions)
  - [Authors](#authors)
  - [License](#license)

### Features

- Stream / Aggregate based event storage and retrieval.
- Global, checkpoint based event retrieval.
<!-- - Flowctl (WIP), a simple command line interface. -->

### Prerequisites

```
docker
```

The easiest way to get up and running is via Docker containers. To use this you need to install Docker:

[Docker installation instructions](https://docs.docker.com/get-docker)

### Installing

Once you've installed Docker, you can execute the following commands to start an EventflowDB instance with a persistent, named volume:

```shell
docker volume create eventflowdb

docker run -d -v eventflowdb:/data -p 6543:6543 -p 16543:16543 -p 176543 ghcr.io/eventflowdb/eventflowdb:0.10.0
```

## Configuration

The following environment variables can be used to alter the configuration:

- `GRPC_PORT`: The port on which the gRPC server listens: Defaults: **6543**
- `HTTP_PORT`: The port on which the HTTP server listens: Defaults: **16543**
- `PROM_PORT`: The port on which the Prometheus HTTP endpoint server listens: Defaults: **26543**

## Usage

EventflowDB is using gRPC with Protobuf as its main method of transport and encoding. The [api.proto](proto/api.proto) file is the source of truth for the API.

The HTTP API is based of the same proto file, but uses the JSON tags for its messages.

API Specification:

```protobuf
service EventStore {
    rpc GetStream(GetStreamRequest) returns (GetStreamResponse) {}
    rpc GetGlobalStream(GetGlobalStreamRequest) returns (GetGlobalStreamResponse) {}
    rpc AppendToStream(AppendToStreamRequest) returns (AppendToStreamResponse) {}
    rpc GetEvent(GetEventRequest) returns (Event) {}
    rpc EventCount(EventCountRequest) returns (EventCountResponse) {}
    rpc StreamCount(StreamCountRequest) returns (StreamCountResponse) {}
    rpc ListStreams(ListStreamsRequest) returns (ListStreamsReponse) {}
    rpc Version(VersionRequest) returns (VersionResponse) {}
}
```

## Example

An simple example project for Golang can be found in the [example](example) folder.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

The features on the roadmap in no particular order:

- Asynchronous replication (with etcd for leader election)
- Downstream message broker connectors (such as Kafka, RabbitMQ)
- Source connector for outbox pattern
- Web UI
- Client libraries for other languages

These may change at any point in the future and will be deleted once they're either implemented or if they won't.

## Contributions

Contributions are most welcome. If you are unsure if a certain feature will benefit the project, please open up an issue.

## Authors

- **Kaj Jagtenberg** - _Initial work_ - [KajJagtenberg](https://github.com/kajjagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details