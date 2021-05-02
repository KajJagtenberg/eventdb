# EventflowDB

EventflowDB is a database designed with Event Sourcing in mind.

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

docker build -t <tag> .

docker volume create eventflowdb

docker run -d -v eventflowdb:/data -p 6543:6543 <tag>
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

* Advanced authentication
* ACL
* Projection Engine
* Asynchronous replication (with Raft for leader election)
* Optional synchronous replication (with Raft)

## Authors

* **Kaj Jagtenberg** - *Initial work* - [KajJagtenberg](https://github.com/KajJagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details
