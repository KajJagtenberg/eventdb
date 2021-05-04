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

Once you've installed Docker, you can execute the following commands to start an EventflowDB instance with a persistent, named volume:

```
docker volume create eventflowdb

docker run -d -v eventflowdb:/data -e PASSWORD=<secure password> -p 6543:6543 docker.pkg.github.com/kajjagtenberg/eventflowdb/eventflowdb:latest
```

Eventhough the "latest" tag is specified in the above command, we recommend you pin the tag to prevent experiencing any unexpected changes.

## Configuration

The following environment variables can be used:

- `PORT`: The port on which the instance used: Defaults: **6543**
- `DATA`: Location of the persisted data (inside the container). Defaults: **/data**
- `PASSWORD`: Clients need to use this password to authenticate to the server. No defaults. We recommend you change this

## Usage

EventflowDB makes use of the [REdis Serialization Protocol](https://redis.io/topics/protocol) which is a simple, easy to understand, text-based protocol. Because of this writing a client for EventflowDB is relatively easy. One can make use of many existing Redis clients available for a large amount of programming languages.

This repository contains a ready to use Golang client which implements all the current features. Support for other languages is on the roadmap, but feel free to contribute your own client libraries.

Commands either use no arguments at all, or a single JSON formatted object which contains the arguments for the given command. Thus commands are either 1 or 2 'words' long, encoded according to the [RESP specification](https://redis.io/topics/protocol).

The following commands are currently supported:

| Command        | Argument                                                        | Description                                                                                                                                                                                                                    | Notes                                                                                                  |
| -------------- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| PING           | -                                                               | Returns PONG                                                                                                                                                                                                                   |                                                                                                        |
| AUTH           | \<password\>                                                    | Authenticates the session                                                                                                                                                                                                      | Only needs to be done after making the connection                                                      |
| CHECKSUM       | -                                                               | Returns the ID of the latest event and the checksum for the entire event log                                                                                                                                                   |                                                                                                        |
| EVENTCOUNT     | -                                                               | Returns the number of events stored                                                                                                                                                                                            | Calling EVENTCOUNTEST is cheaper                                                                       |
| EVENTCOUNTEST  | -                                                               | Returns the number of events stored as of the last periodic index                                                                                                                                                              |                                                                                                        |
| STREAMCOUNT    | -                                                               | Returns the number of streams stored                                                                                                                                                                                           |                                                                                                        |
| STREAMCOUNTEST | -                                                               | Returns the number of streams stored as of the last periodic index                                                                                                                                                             |                                                                                                        |
| ADD            | {"stream":"<uuidv4>","version":<integer>,"events":[\<events\>]} | Appends 1 or more events atomically to a specified stream. If the specified version does not match, a concurrency error will be returned. When successful, the newly stored events will be returned.                           | See the [Event Format](https://github.com/KajJagtenberg/eventflowdb#event-format) for more information |
| GET            | {"stream":"\<uuidv4>","version":\<integer>,"limit":\<integer>}  | Returns events from the given stream, offset by the version, limited by the limit. If the stream contains no events past the given version, then none will be returned.                                                        | Similar to LIMIT and OFFSET in most relational databases                                               |
| GETALL         | {"offset":"\<ulid>","limit":\<integer>}                         | Returns events from all streams. It will skip all events before the vent in the log with the given ID. Returns no more than the amount of events specified by 'limit'                                                          |                                                                                                        |
| QUIT           | Closes the connection                                           |                                                                                                                                                                                                                                |                                                                                                        |
| SIZE           | -                                                               | Returns the size in the form of an array with 2 elements. The first element contains an integer with the size of the database in bytes. The second element contains a string with a human friendly representation of the size. |                                                                                                        |
| UPTIME         | -                                                               | Returns the uptime of the server in human readable format                                                                                                                                                                      |                                                                                                        |
| VERSION        | -                                                               | Returns the version of the database                                                                                                                                                                                            |

### Event Format

The event data must adhere to the following format, otherwise an error will be returned.

```javascript
{
  "type": "AccountOpened", // string
  "data": "{"id": 1, "name": "John Doe"}", // string
  "metadata": "{"user": "3df976f9-f7e6-47e9-a9d0-9e19e451a23e"}", // string. optional
  "causation_id": "0000000000XS4M8WSZ1DW0Z2HT", //ULID in string form, points to the id of the event that caused it. optional
  "correlation_id": "0000000000XS4M8WSZ1DW0Z2HT", //ULID in string form, points to the id of the original event that set the reaction in motion. optional
}
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

The features on the roadmap in no particular order:

- Advanced authentication
- TLS support
- ACL
- Choosable argument encodings (msgpack, protobuf)
- Projection Engine
- Asynchronous replication (with Raft for leader election)
- Optional synchronous replication (with Raft)
- HTTP API
- Backups
- Pub/Sub notifications
- Downstream message broker connectors (such as Kafka, RabbitMQ)
- Web UI / Terminal UI
- Client libraries for other languages

## Contributions

Contributions are most welcome. If you are unsure if a certain feature will benefit the project, please open up an issue.

Discussions are good, arguments are not. Be civil.

## Authors

- **Kaj Jagtenberg** - _Initial work_ - [KajJagtenberg](https://github.com/KajJagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details
