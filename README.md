# Cherry

A concurrent, persistent key-value store built from scratch in Go.

Cherry keeps active data in memory, uses a JSONL write-ahead log (WAL) for durability, replays the WAL on startup for recovery, and exposes the store through an HTTP API and CLI client.

## Features

* In-memory key-value storage using `map[string]string`
* Concurrent access protected by `sync.RWMutex`
* JSONL write-ahead log for persistent writes
* WAL is synced before the in-memory state is updated
* Startup recovery by replaying WAL records in order
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

Stop the server, start it again:

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
Sync WAL to disk
  ↓
Update in-memory map
  ↓
Return success
```

The in-memory state is updated only after the WAL append and sync succeed.

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

Cherry appends and syncs a WAL record before changing the in-memory state. This keeps the persisted operation ahead of the in-memory state.

### `RWMutex`

A single `RWMutex` keeps the concurrency model simple and easy to reason about. More advanced approaches such as sharding or lock striping are outside the scope of this version.

### Synchronous WAL writes

Each write is synced before being treated as successful. This favors durability over write throughput.

### JSONL WAL

JSON Lines keeps the log simple, readable, and easy to replay. A more sophisticated log format could improve efficiency but would add complexity.

## Limitations

Cherry is intentionally a small single-node project.

* Active state is kept entirely in memory
* WAL growth is currently unbounded
* No snapshots or WAL compaction
* No replication or clustering
* No recovery strategy for truncated or partially written WAL records

There is also an edge case around filesystem failures: a WAL record may be written before the subsequent `Sync` operation fails. In that case, the operation can be reported as unsuccessful even though the record may remain in the WAL. This version does not attempt transactional rollback of filesystem state.

## License

MIT License
