<p align="center">
<img src="https://github.com/andygeiss/cloud-native-utils/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Utils

[![Go Reference](https://pkg.go.dev/badge/github.com/andygeiss/cloud-native-utils.svg)](https://pkg.go.dev/github.com/andygeiss/cloud-native-utils)
[![License](https://img.shields.io/github/license/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/releases)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

A modular Go library providing reusable utilities for building cloud-native applications.

---

```bash
go get github.com/andygeiss/cloud-native-utils
```

```go
// Wrap any function with a circuit breaker, then a retry.
fn := stability.Breaker(callPaymentAPI, 3)
fn = stability.Retry(fn, 5, time.Second)
out, err := fn(ctx, in)
```

**Requirements:** Go 1.27 or later.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Running Tests](#running-tests)
- [Baseline deviations](#baseline-deviations)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Cloud Native Utils is a collection of small, focused Go packages designed to be imported independently. There is no monolithic framework—each package addresses a single concern and can be used standalone.

The library covers common cloud-native needs: resilience patterns, structured logging, message dispatching, generic CRUD persistence, security primitives, and HTTP middleware.

---

## Features

| Package | Description |
|---------|-------------|
| **assert** | Minimal test assertion helper (`assert.That`) |
| **consistency** | Transactional event log with JSON file persistence |
| **efficiency** | Channel helpers (`Generate`, `Merge`, `Split`, `Process`), gzip middleware, similarity search (Cosine, Jaccard), sparse data structures (`KeyedSparseSet`, `SparseSharding`) |
| **env** | Generic environment variable parsing (`env.Get[T]`) |
| **event** | Domain event interfaces (`Event`, `EventPublisher`, `EventSubscriber`) |
| **extensibility** | Dynamic Go plugin loading |
| **logging** | Structured JSON logging via `log/slog` |
| **mcp** | Model Context Protocol server for AI tool integrations (Claude Desktop) |
| **messaging** | Publish-subscribe dispatcher (in-memory or Kafka-backed) |
| **resource** | Generic CRUD interface with multiple backends (memory/sharded-sparse/JSON/YAML/SQLite/PostgreSQL) |
| **security** | AES-GCM encryption, password hashing, HMAC, key generation |
| **service** | Context helpers, function wrapper, lifecycle management |
| **slices** | Generic slice utilities (`Map`, `Filter`, `Unique`, etc.) |
| **stability** | Resilience wrappers (circuit breaker, retry, throttle, debounce, timeout) |
| **templating** | HTML template engine with `embed.FS` support |
| **web** | HTTP server, client, routing, sessions, OIDC, session & bearer auth middleware |

---

## Installation

```bash
go get github.com/andygeiss/cloud-native-utils
```

**Requirements:** Go 1.27 or later

---

## Usage

Import only the packages you need:

### Assert

```go
import "github.com/andygeiss/cloud-native-utils/assert"

func TestExample(t *testing.T) {
    result := 42
    assert.That(t, "result should be 42", result, 42)
}
```

### Resource (Generic CRUD)

```go
import "github.com/andygeiss/cloud-native-utils/resource"

// In-memory storage
store := resource.NewInMemoryAccess[string, User]()

// High-performance sharded storage (3-4x faster under concurrency)
store := resource.NewShardedSparseAccess[string, User](32) // 32 shards

// JSON file storage
store := resource.NewJsonFileAccess[string, User]("users.json")

// PostgreSQL storage (requires *sql.DB connection)
store := resource.NewPostgresAccess[string, User](db)
_ = store.Init(ctx) // Creates kv_store table and index

// CRUD operations (same API for all backends)
_ = store.Create(ctx, "user-1", user)
userPtr, _ := store.Read(ctx, "user-1")
_ = store.Update(ctx, "user-1", updatedUser)
_ = store.Delete(ctx, "user-1")
```

### Similarity Search

```go
import (
    "github.com/andygeiss/cloud-native-utils/efficiency"
    "github.com/andygeiss/cloud-native-utils/resource"
)

// Document with sparse vector data
type Document struct {
    Indices []int     // Sorted term indices
    Values  []float64 // TF-IDF values (for cosine)
    Norm    float64   // Pre-computed L2 norm (for cosine)
}

// Create store and populate with documents
store := resource.NewShardedSparseAccess[string, Document](32)
_ = store.Create(ctx, "doc-1", doc1)

// Find similar documents using cosine similarity
results := store.SearchSimilar(ctx, func(doc Document) float64 {
    return efficiency.CosineSimilarity(
        query.Indices, doc.Indices,
        query.Values, doc.Values,
        query.Norm, doc.Norm,
    )
}, resource.SearchOptions{TopK: 10, Threshold: 0.5})

// Find similar documents using Jaccard similarity (for tag sets)
results := store.SearchSimilar(ctx, func(doc Document) float64 {
    return efficiency.JaccardSimilarity(query.Indices, doc.Indices)
}, resource.SearchOptions{TopK: 10})
```

### Stability (Resilience Patterns)

```go
import "github.com/andygeiss/cloud-native-utils/stability"

// Circuit breaker - opens after 3 failures
fn := stability.Breaker(yourFunc, 3)

// Retry with 5 attempts
fn := stability.Retry(yourFunc, 5, time.Second)

// Throttle concurrent executions
fn := stability.Throttle(yourFunc, 10)

// Timeout execution
fn := stability.Timeout(yourFunc, 5*time.Second)
```

### Logging

```go
import "github.com/andygeiss/cloud-native-utils/logging"

logger := logging.NewJsonLogger()
```

### Messaging

```go
import "github.com/andygeiss/cloud-native-utils/messaging"

dispatcher := messaging.NewInternalDispatcher()
_ = dispatcher.Subscribe(ctx, "user.created", handlerFunc)
_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", payload))
```

For Kafka-backed messaging, use `messaging.NewExternalDispatcher()` with `KAFKA_BROKERS` environment variable.

### Event (Domain Events)

```go
import "github.com/andygeiss/cloud-native-utils/event"

// Define a domain event
type UserCreated struct {
    UserID string
}

func (e UserCreated) Topic() string { return "user.created" }

// Use with EventPublisher and EventSubscriber interfaces
var publisher event.EventPublisher = yourPublisher
_ = publisher.Publish(ctx, UserCreated{UserID: "123"})

var subscriber event.EventSubscriber = yourSubscriber
factory := func() event.Event { return &UserCreated{} }
handler := func(e event.Event) error { /* handle event */ return nil }
_ = subscriber.Subscribe(ctx, "user.created", factory, handler)
```

### Env (Environment Variables)

```go
import "github.com/andygeiss/cloud-native-utils/env"

// Generic environment variable parsing with defaults
timeout := env.Get("SERVER_TIMEOUT", 5*time.Second)
maxRetries := env.Get("MAX_RETRIES", 3)
debug := env.Get("DEBUG", false)
rate := env.Get("RATE_LIMIT", 1.5)
name := env.Get("APP_NAME", "my-app")
```

Supported types: `bool`, `int`, `float64`, `string`, `time.Duration`

### Security

```go
import "github.com/andygeiss/cloud-native-utils/security"

// AES-GCM encryption
key := security.GenerateKey()
ciphertext := security.Encrypt([]byte("secret"), key)
plaintext, _ := security.Decrypt(ciphertext, key)

// Password hashing
hash, _ := security.Password([]byte("p@ssw0rd"))
ok := security.IsPasswordValid(hash, []byte("p@ssw0rd"))
```

### Service (Context & Lifecycle)

```go
import "github.com/andygeiss/cloud-native-utils/service"

ctx, cancel := service.Context()
defer cancel()

service.RegisterOnContextDone(ctx, func() {
    // Cleanup logic
})
```

### MCP (Model Context Protocol Server)

```go
import (
    "github.com/andygeiss/cloud-native-utils/mcp"
    "github.com/andygeiss/cloud-native-utils/service"
)

// Create MCP server
server := mcp.NewServer("my-tools", "1.0.0")

// Define tool schema
schema := mcp.NewObjectSchema(
    map[string]mcp.Property{
        "name": mcp.NewStringProperty("Name to greet"),
    },
    []string{"name"},
)

// Register tool with handler
handler := func(ctx context.Context, params mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
    name, _ := params.Arguments["name"].(string)
    return mcp.ToolsCallResult{
        Content: []mcp.ContentBlock{mcp.NewTextContent(fmt.Sprintf("Hello, %s!", name))},
    }, nil
}
server.RegisterTool(mcp.NewTool("greet", "Greets by name", schema, handler))

// Start serving (STDIO transport for Claude Desktop)
ctx, cancel := service.Context()
defer cancel()
server.Serve(ctx)
```

### Web (HTTP Server & Client)

```go
import "github.com/andygeiss/cloud-native-utils/web"

// Create HTTP server with secure defaults
mux := http.NewServeMux()
server := web.NewServer(mux)
server.ListenAndServe()

// Create HTTP client with timeout
client := web.NewClient()

// Create mTLS client
client := web.NewClientWithTLS(certFile, keyFile, caFile)

// Create mux with OIDC, health, liveness, readiness endpoints
//go:embed assets
var efs embed.FS
mux, sessions, idp := web.NewServeMux(ctx, efs)

// Session-based authentication middleware (for web UI)
mux.HandleFunc("GET /protected", web.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value(web.ContextEmail).(string)
    // Handle authenticated request
}))

// Bearer token authentication middleware (for MCP/API endpoints)
// Returns JSON-RPC 2.0 errors on auth failure
verifier := idp.Verifier() // After OIDC provider initialized
mux.HandleFunc("POST /mcp", web.WithBearerAuth(verifier, func(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value(web.ContextEmail).(string)
    subject := r.Context().Value(web.ContextSubject).(string)
    // Handle authenticated MCP request
}))
```

---

## Project Structure

```
cloud-native-utils/
├── assert/          # Test assertions
├── consistency/     # Event logging
├── efficiency/      # Channel helpers, compression, sparse data structures
├── env/             # Environment variable parsing
├── event/           # Domain event interfaces
├── extensibility/   # Plugin loading
├── logging/         # Structured logging
├── mcp/             # MCP server for AI tools
├── messaging/       # Pub-sub dispatchers
├── resource/        # CRUD backends
├── security/        # Cryptographic primitives
├── service/         # Context, lifecycle
├── slices/          # Slice utilities
├── stability/       # Resilience patterns
├── templating/      # Template engine
└── web/             # HTTP server, client, sessions, OIDC
```

For detailed architecture and conventions, see [CLAUDE.md](CLAUDE.md); for the job this library does, see [SPEC.md](SPEC.md).

---

## Running Tests

`make` runs every gate in one go: format, vet, fix, staticcheck, govulncheck, tidy,
test, build.

```bash
make          # the gates against your working tree; run before every commit
make ci       # the same gates against the committed tree; run before every push
```

The inner loop and the extras:

```bash
make test              # go test -race -shuffle=on ./...
make benchmark         # the allocation-sensitive packages
make certs             # a local CA and an mTLS pair, written to security/testdata
make test-integration  # the tests behind the integration build tag
```

`make test-integration` needs what the unit tests do not: a Kafka broker for
`messaging`, an OIDC issuer for `web`, and the certificates `make certs` writes.
`mkcert` must already be installed.

---

## Baseline deviations

This library follows [Andy's engineering baseline](https://github.com/andygeiss/baseline).
Where it differs, it says so here.

### Third-party dependencies

The baseline asks a library for **zero** third-party dependencies, and asks that any
exception be justified. There are seven.

| Dependency | Used by | Why the standard library is not enough |
|---|---|---|
| `github.com/coreos/go-oidc/v3` | web | OIDC discovery and ID-token checking. Hand-rolling it means hand-rolling JWKS rotation and JWT validation — the part of a login flow that must not be homemade. |
| `github.com/jackc/pgx/v5` | resource | Postgres driver and pool. On the baseline's approved list. |
| `github.com/segmentio/kafka-go` | messaging | The Kafka wire protocol. There is no standard-library equivalent, and a Kafka-backed dispatcher is why `messaging` exists. |
| `golang.org/x/crypto` | security | argon2 and bcrypt password hashing. On the approved list. |
| `golang.org/x/oauth2` | web | The authorization-code flow `go-oidc` is built on. Taking `go-oidc` means taking this. |
| `gopkg.in/yaml.v3` | resource | YAML parsing. The standard library has none. |
| `modernc.org/sqlite` | resource | SQLite driver in pure Go, so binaries stay CGO-free. On the approved list. |

Four of them — `go-oidc`, `kafka-go`, `oauth2` and `yaml.v3` — are not on the approved
list in the baseline's `stack/go.md`. Each stays inside the one package named above, so
a project that imports only `slices` or `stability` never builds it.

### Rules met by a different route

These are not waivers. The rule holds; this library reaches it another way.

- **A library must not log.** No package here writes a log line by itself. `logging`
  hands you a configured `*slog.Logger`, and `web.WithLogging` is a middleware you
  install with your own logger. What is worth logging stays your decision.
- **Implementation detail belongs under `internal/`.** There is no `internal/`. Every
  package is public surface on purpose: the module is a set of small packages, not one
  package with parts hidden inside it.

---

## Contributing

Contributions are welcome:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Follows the conventions in [CLAUDE.md](CLAUDE.md)
- Includes tests (`*_test.go` files)
- Passes `make ci`

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
