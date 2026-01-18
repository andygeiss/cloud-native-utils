<p align="center">
<img src="https://github.com/andygeiss/cloud-native-utils/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Utils

[![Go Reference](https://pkg.go.dev/badge/github.com/andygeiss/cloud-native-utils.svg)](https://pkg.go.dev/github.com/andygeiss/cloud-native-utils)
[![License](https://img.shields.io/github/license/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-utils)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-utils)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A modular Go library providing reusable utilities for building cloud-native applications.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Running Tests](#running-tests)
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
| **efficiency** | Channel helpers (`Generate`, `Merge`, `Split`, `Process`) and gzip middleware |
| **env** | Generic environment variable parsing (`env.Get[T]`) |
| **event** | Domain event interfaces (`Event`, `EventPublisher`, `EventSubscriber`) |
| **extensibility** | Dynamic Go plugin loading |
| **logging** | Structured JSON logging via `log/slog` |
| **mcp** | Model Context Protocol server for AI tool integrations (Claude Desktop) |
| **messaging** | Publish-subscribe dispatcher (in-memory or Kafka-backed) |
| **resource** | Generic CRUD interface with multiple backends (memory/JSON/YAML/SQLite/PostgreSQL) |
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

**Requirements:** Go 1.25.4 or later

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

// JSON file storage
store := resource.NewJsonFileAccess[string, User]("users.json")

// PostgreSQL storage (requires *sql.DB connection)
store := resource.NewPostgresAccess[string, User](db)
_ = store.Init(ctx) // Creates kv_store table and index

// CRUD operations
_ = store.Create(ctx, "user-1", user)
userPtr, _ := store.Read(ctx, "user-1")
_ = store.Update(ctx, "user-1", updatedUser)
_ = store.Delete(ctx, "user-1")
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
mux, sessions := web.NewServeMux(ctx, efs)

// Session-based authentication middleware (for web UI)
mux.HandleFunc("GET /protected", web.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value(web.ContextEmail).(string)
    // Handle authenticated request
}))

// Bearer token authentication middleware (for MCP/API endpoints)
// Returns JSON-RPC 2.0 errors on auth failure
verifier := web.IdentityProvider.Verifier() // After OIDC provider initialized
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
├── efficiency/      # Channel helpers, compression
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

For detailed architecture and conventions, see [CONTEXT.md](CONTEXT.md).

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbose
go test -v ./...

# Using just (recommended)
just test
```

---

## Linting

This project uses [golangci-lint](https://golangci-lint.run/) for code quality checks.

```bash
# Run linter
just lint

# Or directly
golangci-lint run ./...
```

Configuration is in `.golangci.yml`.

---

## Contributing

Contributions are welcome:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Follows the conventions in [CONTEXT.md](CONTEXT.md)
- Includes tests (`*_test.go` files)
- Passes `just test` and `just lint`

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
