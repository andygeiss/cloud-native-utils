<p align="center">
<img src="https://github.com/andygeiss/cloud-native-utils/blob/main/logo.png?raw=true" />
</p>

# Cloud Native Utils

[![Go Reference](https://pkg.go.dev/badge/badge/andygeiss/cloud-native-utils.svg)](https://pkg.go.dev/badge/andygeiss/cloud-native-utils)
[![License](https://img.shields.io/github/license/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/cloud-native-utils)](https://github.com/andygeiss/cloud-native-utils/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/cloud-native-utils)](https://goreportcard.com/report/github.com/andygeiss/cloud-native-utils)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/b4e3a9c4859b47f1bc43613970ec8d12)](https://app.codacy.com/gh/andygeiss/cloud-native-utils/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

A collection of high-performance, modular utilities for building cloud-native Go applications.
This library provides battle-tested patterns for testing, data consistency, concurrency,
security, messaging, and service stability.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Assert](#assert)
  - [Consistency](#consistency)
  - [Efficiency](#efficiency)
  - [Extensibility](#extensibility)
  - [i18n](#i18n)
  - [Logging](#logging)
  - [Messaging](#messaging)
  - [Redirecting](#redirecting)
  - [Resource](#resource)
  - [Security](#security)
  - [Service](#service)
  - [Stability](#stability)
  - [Templating](#templating)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Assert** - Simple and expressive test assertions with deep equality checks
- **Consistency** - Transactional log management with event abstractions and file-based persistence
- **Efficiency** - Channel utilities for generating, merging, splitting streams, and sharded key-value stores
- **Extensibility** - Dynamic Go plugin loading for on-the-fly feature integration
- **i18n** - YAML-based internationalization with embedded filesystem support
- **Logging** - Structured JSON logging with HTTP middleware support
- **Messaging** - Publish-subscribe patterns for decoupling local and remote services
- **Redirecting** - HTMX-compatible HTTP redirects for state-changing requests
- **Resource** - Generic CRUD interface with in-memory, JSON, YAML, and SQLite backends
- **Security** - AES-GCM encryption, secure key generation, HMAC hashing, bcrypt passwords, OAuth2/OIDC, and secure HTTP servers
- **Service** - Context-aware function wrappers with lifecycle and signal handling
- **Stability** - Circuit breakers, retries, throttling, debounce, and timeouts
- **Templating** - Template engine with embedded filesystem support

---

## Installation

```bash
go get github.com/andygeiss/cloud-native-utils
```

**Requirements:**
- Go 1.21 or later

---

## Usage

### Assert

A utility function to assert value equality in tests with clear error messages.

```go
import (
    "testing"
    "github.com/andygeiss/cloud-native-utils/assert"
)

func TestExample(t *testing.T) {
    result := 42
    assert.That(t, "result should be 42", result, 42)
}
```

### Consistency

Transactional log management with events and JSON file-based persistence.

```go
import "github.com/andygeiss/cloud-native-utils/consistency"

// Create events for logging state changes
event := consistency.NewEvent(consistency.EventTypeCreate, "user", "123", userData)

// Use JsonFileLogger for persistent event storage
logger := consistency.NewJsonFileLogger("events.json")
logger.Log(event)
```

### Efficiency

Utilities for concurrent stream processing and data partitioning.

```go
import "github.com/andygeiss/cloud-native-utils/efficiency"

// Generate a read-only channel from values
ch := efficiency.Generate(ctx, 1, 2, 3, 4, 5)

// Merge multiple channels into one
merged := efficiency.Merge(ctx, ch1, ch2, ch3)

// Split a channel into multiple outputs
outputs := efficiency.Split(ctx, input, 3)

// Process items concurrently
efficiency.Process(ctx, input, func(item int) {
    // Handle each item
})
```

### Extensibility

Dynamically load external Go plugins at runtime.

```go
import "github.com/andygeiss/cloud-native-utils/extensibility"

// Load a symbol from a plugin file
symbol, err := extensibility.LoadPlugin("./plugins/myplugin.so", "MyFunction")
```

### i18n

YAML-based internationalization with embedded filesystem support.

```go
import "github.com/andygeiss/cloud-native-utils/i18n"

//go:embed translations/*.yaml
var translationsFS embed.FS

translations := i18n.NewTranslations()
translations.Load(translationsFS, "en", "translations/en.yaml")
translations.Load(translationsFS, "de", "translations/de.yaml")

// Get translated text
text := translations.Get("en", "greeting.hello")
```

### Logging

Structured JSON logging with HTTP middleware.

```go
import "github.com/andygeiss/cloud-native-utils/logging"

// Create a JSON logger
logger := logging.NewJsonLogger(os.Stdout)

// Use middleware for HTTP request logging
handler := logging.Middleware(logger)(yourHandler)
```

### Messaging

Publish-subscribe patterns for decoupling services.

```go
import "github.com/andygeiss/cloud-native-utils/messaging"

// Create a dispatcher for internal messaging
dispatcher := messaging.NewInternalDispatcher()

// Subscribe to a topic
dispatcher.Subscribe("user.created", func(msg messaging.Message) {
    // Handle message
})

// Publish a message
dispatcher.Publish("user.created", payload)
```

### Redirecting

HTMX-compatible HTTP redirects for POST/PUT/DELETE requests.

```go
import "github.com/andygeiss/cloud-native-utils/redirecting"

// Wrap handlers to redirect state-changing requests to GET endpoints
handler := redirecting.Middleware("/success")(yourHandler)
```

### Resource

Generic CRUD interface with multiple backend implementations.

```go
import "github.com/andygeiss/cloud-native-utils/resource"

// In-memory storage
store := resource.NewInMemoryAccess[string, User]()

// JSON file storage
store := resource.NewJsonFileAccess[string, User]("users.json")

// SQLite storage
store := resource.NewSqliteAccess[string, User](db, "users")

// CRUD operations
store.Create(ctx, "user-1", user)
user, err := store.Read(ctx, "user-1")
store.Update(ctx, "user-1", updatedUser)
store.Delete(ctx, "user-1")
users, err := store.List(ctx)
```

### Security

Comprehensive security utilities for cloud-native applications.

```go
import "github.com/andygeiss/cloud-native-utils/security"

// AES-GCM encryption/decryption
key := security.GenerateKey()
encrypted := security.Encrypt(key, plaintext)
decrypted := security.Decrypt(key, encrypted)

// Password hashing with bcrypt
hash := security.HashPassword(password)
valid := security.VerifyPassword(hash, password)

// Generate secure IDs
id := security.GenerateID()

// PKCE for OAuth2
verifier, challenge := security.GeneratePKCE()

// Secure HTTP server with health probes
server := security.NewServer(":8443", handler)
```

### Service

Context-aware function wrappers with lifecycle management.

```go
import "github.com/andygeiss/cloud-native-utils/service"

// Wrap a function with context support
fn := service.Wrap(func(ctx context.Context) error {
    // Your service logic
    return nil
})

// Register cleanup on context cancellation
service.RegisterOnContextDone(ctx, func() {
    // Cleanup logic
})
```

### Stability

Patterns for building resilient services.

```go
import "github.com/andygeiss/cloud-native-utils/stability"

// Circuit breaker - opens after threshold failures
fn := stability.Breaker(yourFunc, 3)

// Retry with configurable attempts
fn := stability.Retry(yourFunc, 5, time.Second)

// Throttle to limit concurrent executions
fn := stability.Throttle(yourFunc, 10)

// Debounce to delay execution
fn := stability.Debounce(yourFunc, 500*time.Millisecond)

// Timeout to limit execution time
fn := stability.Timeout(yourFunc, 5*time.Second)
```

### Templating

Template engine with embedded filesystem support.

```go
import "github.com/andygeiss/cloud-native-utils/templating"

//go:embed templates/*.html
var templatesFS embed.FS

engine := templating.NewEngine(templatesFS)
engine.Parse("templates/*.html")
engine.Render(w, "page.html", data)
```

---

## Technologies Used

- **Go** (1.25+) - Primary programming language
- **AES-GCM** - Authenticated encryption
- **bcrypt** - Password hashing
- **OAuth2/OIDC** - Authentication protocols
- **SQLite** - Embedded database support
- **Kafka** - Message queue integration
- **YAML** - Configuration and i18n files

---

## Contributing

Contributions are welcome! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

Please ensure your code:
- Follows Go best practices and idioms
- Includes tests for new functionality
- Passes all existing tests (`go test ./...`)
- Has appropriate documentation

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
