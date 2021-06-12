# EventflowDB

EventflowDB is a database designed with Event Sourcing in mind.

The current version is subject to change and the API may break at any time. Be advised.

- [EventflowDB](#eventflowdb)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installing](#installing)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Event Format](#event-format)
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

docker run -d -v eventflowdb:/data -e PASSWORD=<secure password> -p 6543:6543 -p 16543:16543 kajjagtenberg/eventflowdb:0.7.0
```

## Configuration

The following environment variables can be used to alter the configuration:

- `PORT`: The port on which the RESP server listens: Defaults: **6543**
- `DATA`: Location of the persisted data (inside the container). Defaults: **/data**
- `TLS_ENABLED`: true/false. Enable TLS for RESP and HTTP API. Defaults: **false**
- `TLS_CERT_FILE`: Location of the certificate. Defaults: **certs/cert.pem**
- `TLS_KEY_FILE`: Location of the key. Defaults: **certs/key.pem**

## Usage

EventflowDB makes use of the [REdis Serialization Protocol](https://redis.io/topics/protocol) which is a simple, easy to understand, text-based protocol. Because of this writing a client for EventflowDB is relatively easy. One can make use of many existing Redis clients available for a large amount of programming languages.

This repository contains a ready to use Golang client which implements all the current features. Support for other languages is on the roadmap, but feel free to contribute your own client libraries.

Commands either use no arguments at all, or a single JSON formatted object, which contains the arguments for the given command. Thus commands are either 1 or 2 'words' long, encoded according to the [RESP specification](https://redis.io/topics/protocol).

The following commands are currently supported:

| Command        | Shorthand | Argument                                                        | Description                                                                                                                                                                                                                    | Notes                                                                                                  |
| -------------- | --------- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| PING           | P         | -                                                               | Returns PONG                                                                                                                                                                                                                   |                                                                                                        |
| AUTH           | -         | \<password\>                                                    | Authenticates the session                                                                                                                                                                                                      | Only needs to be done after making the connection                                                      |
| CHECKSUM       | CH        | -                                                               | Returns the ID of the latest event and the checksum for the entire event log                                                                                                                                                   |                                                                                                        |
| EVENTCOUNT     | EC        | -                                                               | Returns the number of events stored                                                                                                                                                                                            | Calling EVENTCOUNTEST is cheaper                                                                       |
| EVENTCOUNTEST  | ECE       | -                                                               | Returns the number of events stored as of the last periodic index                                                                                                                                                              |                                                                                                        |
| STREAMCOUNT    | SC        | -                                                               | Returns the number of streams stored                                                                                                                                                                                           | Calling STREAMCOUNTEST is cheaper                                                                      |
| STREAMCOUNTEST | SCE       | -                                                               | Returns the number of streams stored as of the last periodic index                                                                                                                                                             |                                                                                                        |
| ADD            | A         | {"stream":"<uuidv4>","version":<integer>,"events":[\<events\>]} | Appends 1 or more events atomically to a specified stream. If the specified version does not match, a concurrency error will be returned. When successful, the newly stored events will be returned.                           | See the [Event Format](https://github.com/kajjagtenberg/eventflowdb#event-format) for more information |
| GET            | G         | {"stream":"\<uuidv4>","version":\<integer>,"limit":\<integer>}  | Returns events from the given stream, offset by the version, limited by the limit. If the stream contains no events past the given version, then none will be returned.                                                        | Similar to LIMIT and OFFSET in most relational databases                                               |
| GETALL         | GA        | {"offset":"\<ulid>","limit":\<integer>}                         | Returns events from all streams. It will skip all events before the vent in the log with the given ID. Returns no more than the amount of events specified by 'limit'                                                          |                                                                                                        |
| SIZE           | S         | -                                                               | Returns the size in the form of an array with 2 elements. The first element contains an integer with the size of the database in bytes. The second element contains a string with a human friendly representation of the size. |                                                                                                        |
| UPTIME         | UP        | -                                                               | Returns the uptime of the server in human readable format                                                                                                                                                                      |                                                                                                        |
| VERSION        | V         | -                                                               | Returns the version of the database                                                                                                                                                                                            |

Access to the API is also available via HTTP. All commands are done via POST request with the following url:

```
POST http://<hostname>:16543/api/:cmd
```

Both the full name and the shorthand name of commands can be used to send the command. The body the JSON arguments if neccessary.

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

All methods that return a list of events will be in the following format:

```javascript
{
  "id": "0000000000XS4M8WSZ1DW0Z2HT",
  "stream": "7c682570-fce1-4518-9ea3-486648301183",
  "version": 0,
  "type": "AccountOpened",
  "data": "{"id": 1, "name": "John Doe"}",
  "metadata": "{"user": "3df976f9-f7e6-47e9-a9d0-9e19e451a23e"}",
  "causation_id": "0000000000XS4M8WSZ1DW0Z2HT",
  "correlation_id": "0000000000XS4M8WSZ1DW0Z2HT",
  "added_at": "2021-05-05T11:02:56.372078255Z"
}
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kajjagtenberg/eventflowdb/tags).

## Roadmap

The features on the roadmap in no particular order:

- Advanced authentication
- TLS support
- ACL or other authorization scheme
- Choosable argument encodings (msgpack, protobuf)
- Projection Engine
- Asynchronous replication (with external leader election)
- Optional synchronous replication (with Raft)
- Backups
- Pub/Sub notifications
- Downstream message broker connectors (such as Kafka, RabbitMQ)
- Web UI
- Client libraries for other languages

## Contributions

Contributions are most welcome. If you are unsure if a certain feature will benefit the project, please open up an issue.

## Authors

- **Kaj Jagtenberg** - _Initial work_ - [KajJagtenberg](https://github.com/kajjagtenberg)

See also the list of [contributors](https://github.com/kajjagtenberg/eventflowdb/contributors) who participated in this project.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details
