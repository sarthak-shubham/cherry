# Cherry

A concurrent, persistent key-value store built from scratch in Go.

Cherry keeps active data in memory, uses a JSONL write-ahead log (WAL) for persistence, replays the WAL on startup for recovery, and exposes the store through an HTTP API and CLI client.

## Features

* In-memory key-value storage using `map[string]string`
* Concurrent access protected by `sync.RWMutex`
* JSONL write-ahead log for persistent writes
* WAL write validation and rollback on failed or partial writes
* Startup recovery by replaying WAL records in order
* Automatic removal of an incomplete final WAL record after a crash
* HTTP/JSON API for `GET`, `PUT`, and `DELETE`
* CLI client for interacting with the running server
* Unit, integration, API, client, and concurrency tests
* Verified with Go's race detector

## Quick Start

### Requirements

* Go 1.27+

### Start the server

```bash
go run .
```

The server listens on:

```text
http://localhost:8080
```

### Use the CLI

Open another terminal:

```bash
go run ./client set name Robert
go run ./client get name
go run ./client delete name
```

Example:

```text
name = Robert

name = Robert

deleted name
```

### Test persistence

Set some values:

```bash
go run ./client set name Robert
go run ./client set city Tokyo
```

Stop the server and start it again:

```bash
go run .
```

Then retrieve the values:

```bash
go run ./client get name
go run ./client get city
```

The values are restored by replaying the WAL during startup.

## API

| Method   | Endpoint    | Description               |
| -------- | ----------- | ------------------------- |
| `PUT`    | `/kv/{key}` | Create or replace a value |
| `GET`    | `/kv/{key}` | Retrieve a value          |
| `DELETE` | `/kv/{key}` | Delete a key              |

### Set a value

```http
PUT /kv/name
Content-Type: application/json

{"value":"Robert"}
```

### Get a value

```http
GET /kv/name
```

Response:

```json
{
  "key": "name",
  "value": "Robert"
}
```

### Delete a key

```http
DELETE /kv/name
```

## How It Works

### Write path

```text
Client
  ↓
HTTP API
  ↓
Store
  ↓
Write lock
  ↓
Append WAL record
  ↓
Sync WAL
  ↓
Update in-memory map
  ↓
Return success
```

Cherry writes the WAL record before updating the in-memory map.

If the WAL write fails or is incomplete, Cherry attempts to roll the WAL back to its previous size and returns an error for that request. The in-memory map is not changed.

If the WAL record is written successfully but `Sync()` fails, Cherry logs a warning and continues serving. The operation is applied to the in-memory map, but its durability after a crash cannot be guaranteed.

A failure affecting one write therefore does not terminate the server or prevent other requests from being served.

### Read path

Reads are served directly from the in-memory map:

```text
GET
 ↓
RLock
 ↓
Map lookup
 ↓
RUnlock
```

### Recovery

Cherry rebuilds its in-memory state from the WAL when the application starts:

```text
WAL
  ↓
Replay records in order
  ↓
Rebuild map
  ↓
Create Store
  ↓
Start HTTP server
```

A `set` record creates or replaces a value, while a `delete` record removes the key.

Cherry handles different WAL startup conditions differently:

* **WAL does not exist:** Cherry starts with an empty map. This is treated as a fresh start.
* **WAL exists but is empty:** Cherry starts with an empty map.
* **WAL ends with an incomplete final record:** Cherry truncates the incomplete tail and recovers all preceding complete records.
* **WAL cannot be opened or read:** Recovery fails because Cherry cannot safely determine the persisted state.
* **A complete WAL record contains invalid JSON or an unknown operation:** Recovery fails rather than silently discarding persisted data.

Cherry therefore does not treat an existing but unreadable or corrupt WAL as an empty database. Doing so could silently discard previously persisted data.

## Concurrency

The store uses `sync.RWMutex` to protect its in-memory map.

Concurrent reads use a read lock, while writes use an exclusive lock.

The concurrency test runs multiple goroutines performing `GET`, `SET`, and `DELETE` operations against shared keys. The project is also verified with Go's race detector:

```bash
go test -race ./...
```

## Testing

Run all tests:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```

The test suite covers:

* Store operations and error handling
* WAL persistence and recovery
* HTTP API behavior
* CLI request handling
* Persistence and recovery integration
* Concurrent store access

## Design Tradeoffs

### In-memory state

Keeping the active dataset in a Go map keeps the implementation simple and makes reads straightforward, but the working dataset must fit in memory.

### WAL before memory mutation

Cherry writes the WAL record before changing the in-memory state.

A failed or incomplete WAL write prevents the in-memory mutation.

If the record is written successfully but `Sync()` fails, Cherry favors availability and continues serving while warning that durability is uncertain.

### `RWMutex`

A single `RWMutex` keeps the concurrency model simple and easy to reason about. More advanced approaches such as sharding or lock striping are outside the scope of this version.

### Synchronous WAL writes

Cherry attempts to sync each WAL record before treating its durability as guaranteed.

If `Sync()` fails after the record has been successfully written, Cherry logs the failure and continues serving rather than rejecting the operation.

This favors availability over guaranteeing durability in the face of a filesystem sync failure.

### JSONL WAL

JSON Lines keeps the log simple, readable, and easy to replay. A more sophisticated log format could improve efficiency but would add complexity.

## Limitations

Cherry is intentionally a small single-node project.

* Active state is kept entirely in memory
* WAL growth is currently unbounded
* No snapshots or WAL compaction
* No replication or clustering
* No distributed recovery or failover

Filesystem failures can still create durability limitations.

In particular, if a WAL record is successfully written but `Sync()` fails, Cherry accepts the operation in memory while logging a warning. The operation may therefore be lost after a crash if the data was not durably flushed to disk.

More advanced storage and recovery mechanisms are outside the scope of the current version.

## License

MIT License
