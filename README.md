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
| **extensibility** | Dynamic Go plugin loading |
| **imaging** | QR code generation |
| **logging** | Structured JSON logging via `log/slog` with HTTP middleware |
| **messaging** | Publish-subscribe dispatcher (in-memory or Kafka-backed) |
| **redirecting** | HTMX-compatible PRG redirects |
| **resource** | Generic CRUD interface with multiple backends (memory/JSON/YAML/SQLite) |
| **security** | AES-GCM encryption, password hashing, OIDC, HTTP server |
| **service** | Context helpers, function wrapper, lifecycle management |
| **slices** | Generic slice utilities (`Map`, `Filter`, `Unique`, etc.) |
| **stability** | Resilience wrappers (circuit breaker, retry, throttle, debounce, timeout) |
| **templating** | HTML template engine with `embed.FS` support |

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
handler := logging.WithLogging(logger, yourHandler)
```

### Messaging

```go
import "github.com/andygeiss/cloud-native-utils/messaging"

dispatcher := messaging.NewInternalDispatcher()
_ = dispatcher.Subscribe(ctx, "user.created", handlerFunc)
_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", payload))
```

For Kafka-backed messaging, use `messaging.NewExternalDispatcher()` with `KAFKA_BROKERS` environment variable.

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

// Configured HTTP server
server := security.NewServer(mux)
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

---

## Project Structure

```
cloud-native-utils/
├── assert/          # Test assertions
├── consistency/     # Event logging
├── efficiency/      # Channel helpers, compression
├── extensibility/   # Plugin loading
├── imaging/         # QR code generation
├── logging/         # Structured logging
├── messaging/       # Pub-sub dispatchers
├── redirecting/     # HTTP redirects
├── resource/        # CRUD backends
├── security/        # Encryption, auth, server
├── service/         # Context, lifecycle
├── slices/          # Slice utilities
├── stability/       # Resilience patterns
└── templating/      # Template engine
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
```

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
- Passes `go test ./...` and `go vet ./...`

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
