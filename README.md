# EventflowDB

EventflowDB is a database designed with Event Sourcing in mind.

The current version is subject to change and the API may break at any time. Be advised.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

```
docker
```

The easiest way to get up and running is via Docker containers. To use this you need to install Docker:

[Docker installation instructions](https://docs.docker.com/get-docker)

### Installing

Once you've installed Docker, you can execute the following commands to start an EventflowDB instance:

```
git clone https://github.com/KajJagtenberg/eventflowdb eventflowdb

cd eventflowdb

docker build -t eventflowdb:latest .

docker volume create eventflowdb

docker run -d -v eventflowdb:/data -e PASSWORD=<secure password> -p 6543:6543 eventflowdb:latest
```

### Configuration

The following environment variables can be used:

* `PORT`: The port on which the instance used: Defaults: __6543__
* `LANG`: System language. Defaults: __en_US.UTF-8__
* `DATA`: Location of the persisted data (inside the container). Defaults: __/data__
* `PASSWORD`: Clients need to use this password to authenticate to the server. No defaults. We recommend you change this

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

- Advanced authentication
- TLS support
- ACL
- Choosable payload encodings (msgpack, protobuf)
- Projection Engine
- Asynchronous replication (with Raft for leader election)
- Optional synchronous replication (with Raft)
- HTTP API
- Backups
- Web UI
- Pub/Sub notifications
- Downstream message broker connectors (such as Kafka, RabbitMQ)

## Authors

- **Kaj Jagtenberg** - _Initial work_ - [KajJagtenberg](https://github.com/KajJagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details
