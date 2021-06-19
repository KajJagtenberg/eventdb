# EventflowDB

EventflowDB is a database designed with Event Sourcing in mind.

- [EventflowDB](#eventflowdb)
  - [Getting Started](#getting-started)
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

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

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

docker run -d -v eventflowdb:/data -p 6543:6543 ghcr.io/eventflowdb:latest
```

## Configuration

The following environment variables can be used to alter the configuration:

- `GRPC_PORT`: The port on which the gRPC server listens: Defaults: **6543**
- `DATA`: Location of the persisted data (inside the container). Defaults: **/data**
- `TLS_ENABLED`: true/false. Enable TLS for RESP and HTTP API. Defaults: **false**
- `TLS_CERT_FILE`: Location of the certificate. Defaults: **certs/cert.pem**
- `TLS_KEY_FILE`: Location of the key. Defaults: **certs/key.pem**
- `IN_MEMORY`: Whether the data should resize in memory only. Defaults: **false**

## Usage

EventflowDB is using gRPC with Protobuf as its method of transport and encoding. The [api.proto](proto/api.proto) file is the source of truth for the API.

## Example

An example for Golang can be found in the [example](example) folder.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

The features on the roadmap in no particular order:

- TLS Client Authentication
- ACL or other authorization scheme
- Projection Engine
- Asynchronous replication (with etcd for leader election)
- Downstream message broker connectors (such as Kafka, RabbitMQ)
- Web UI
- Client libraries for other languages
- CLI
- Prometheus metrics

## Contributions

Contributions are most welcome. If you are unsure if a certain feature will benefit the project, please open up an issue.

## Authors

- **Kaj Jagtenberg** - _Initial work_ - [KajJagtenberg](https://github.com/kajjagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details
