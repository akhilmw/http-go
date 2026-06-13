# HTTP-Go

A lightweight HTTP/1.1 server built from scratch in Go using raw TCP sockets,
without relying on Go's `net/http` server. This project was built to understand
the internals of HTTP, request parsing, response generation, chunked transfer
encoding, and streaming.

## Features

### Request Parsing

- Parse HTTP request lines
- Parse headers with case-insensitive lookup
- Support duplicate headers
- Parse request bodies using `Content-Length`
- Stateful request parsing

### Response Generation

- HTTP status line generation
- Header serialization
- Automatic `Content-Length`
- Custom response writer abstraction
- Plain text and HTML responses

### Server Architecture

- Concurrent connection handling with goroutines
- Graceful shutdown
- Handler-based API inspired by Go's `net/http`

### Chunked Transfer Encoding

- Streaming responses with `Transfer-Encoding: chunked`
- Hexadecimal chunk sizes
- Chunk terminators
- HTTP trailers

### Proxy Support

Requests to:

```text
/httpbin/<path>
```

are proxied to:

```text
https://httpbin.org/<path>
```

with support for streaming chunked responses and trailers.

## Project Structure

```text
httpgo/
├── httpserver/
│   └── main.go
│
├── internal/
│   ├── headers/
│   │   ├── headers.go
│   │   └── headers_test.go
│   │
│   ├── request/
│   │   ├── request.go
│   │   └── request_test.go
│   │
│   ├── response/
│   │   ├── response.go
│   │   └── response_test.go
│   │
│   └── server/
│       └── server.go
│
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── messages.txt
```

## Architecture

```text
TCP Connection
       ↓
RequestFromReader()
       ↓
request.Request
       ↓
handler(*response.Writer, req)
       ↓
Response Writer
├── Status Line
├── Headers
├── Body
├── Chunked Encoding
└── Trailers
       ↓
Client
```

## Running

Start the server:

```bash
go run ./httpserver
```

The server listens on port **42069**.

## Examples

### Basic Request

```bash
curl localhost:42069
```

### Error Routes

```bash
curl localhost:42069/yourproblem
curl localhost:42069/myproblem
```

### Proxy Requests

```bash
curl localhost:42069/httpbin/get
curl localhost:42069/httpbin/html
```

### View Raw Chunked Responses

```bash
curl --raw -i localhost:42069/httpbin/html
```

or:

```bash
echo -e "GET /httpbin/stream/3 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069
```

## Testing

Run all tests:

```bash
go test ./...
```

## Concepts Explored

- TCP sockets
- HTTP/1.1 message format
- Request parsing
- Header parsing
- Body parsing
- Response generation
- Status codes
- `Content-Length`
- Chunked transfer encoding
- HTTP trailers
- Streaming
- Reverse proxying
- Goroutines
- Interfaces (`io.Reader`, `io.Writer`)
- State machines

## Inspiration

Built while following Boot.dev's **Learn the HTTP Protocol in Go** course with
the goal of understanding what lies beneath:

```go
http.ListenAndServe(":8080", nil)
```

by implementing the underlying pieces from scratch.

## Future Improvements

- HTTP keep-alive connections
- Router with path parameters
- Middleware support
- Logging
- TLS/HTTPS
- Compression
- WebSockets
- HTTP/2 support
- Static file serving
- Request timeouts
- Connection pooling
