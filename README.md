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

## **Installation**

```bash
go get github.com/andygeiss/cloud-native-utils
```

---

## **Module Features**

### [`assert`](assert/) — Testing Utilities

Simplifies test assertions with a clean, consistent API.

| Function | Description |
|----------|-------------|
| `That(t, desc, got, expected)` | Asserts that two values are deeply equal |

```go
import "github.com/andygeiss/cloud-native-utils/assert"

func Test_Example(t *testing.T) {
    result := Add(2, 3)
    assert.That(t, "sum must equal 5", result, 5)
}
```

---

### [`consistency`](consistency/) — Transactional Log Management

Implements write-ahead logging (WAL) patterns for data consistency with event sourcing support.

| Type | Description |
|------|-------------|
| `Event[K, V]` | Generic event with key, value, sequence number, and event type |
| `EventType` | Enum for `EventTypeDelete` and `EventTypePut` operations |
| `JsonFileLogger[K, V]` | File-based logger with JSON persistence and async writes |
| `Logger[K, V]` | Interface for transactional logging |

```go
import "github.com/andygeiss/cloud-native-utils/consistency"

logger := consistency.NewJsonFileLogger[string, Task]("tasks.log")
defer logger.Close()

logger.Log(consistency.Event[string, Task]{
    EventType: consistency.EventTypePut,
    Key:       "task-1",
    Value:     task,
})
```

---

### [`efficiency`](efficiency/) — Concurrency & Data Structures

Utilities for concurrent processing, channel operations, and high-performance data structures.

| Function/Type | Description |
|---------------|-------------|
| `Generate[T](values...)` | Creates a read-only channel from variadic values |
| `Merge[T](channels...)` | Combines multiple channels into one |
| `Split[T](in, num)` | Splits one channel into multiple channels |
| `Process[IN, OUT](in, fn)` | Concurrent processing with CPU-based worker pool |
| `Sharding[K, V]` | Thread-safe sharded key-value store |
| `SparseSet[T]` | Memory-efficient sparse set data structure |
| `WithCompression(handler)` | HTTP middleware for gzip compression |

```go
import "github.com/andygeiss/cloud-native-utils/efficiency"

// Channel pipeline
ch := efficiency.Generate(1, 2, 3, 4, 5)
results, errors := efficiency.Process(ch, processFunc)

// Sharded storage for high-concurrency scenarios
shards := efficiency.NewSharding[string, User](16)
shards.Put("user-1", user)
user, exists := shards.Get("user-1")
```

---

### [`extensibility`](extensibility/) — Plugin System

Dynamically loads external Go plugins at runtime for hot-swappable features.

| Function | Description |
|----------|-------------|
| `LoadPlugin[T](path, symbolName)` | Loads a typed symbol from a `.so` plugin file |

```go
import "github.com/andygeiss/cloud-native-utils/extensibility"

fn, err := extensibility.LoadPlugin[func(int) int]("./plugins/math.so", "Double")
result := fn(21) // 42
```

---

### [`i18n`](i18n/) — Internationalization

YAML-based translation management with dot-notation key lookup.

| Type/Method | Description |
|-------------|-------------|
| `Translations` | Thread-safe translation store |
| `Load(efs, lang, path)` | Loads translations from embedded YAML files |
| `T(lang, key)` | Returns translated string for a dot-separated key |
| `TMap(lang, keys...)` | Returns a map of translations for template rendering |

```go
import "github.com/andygeiss/cloud-native-utils/i18n"

translations := i18n.NewTranslations()
translations.Load(efs, "en", "assets/i18n/en.yaml")
translations.Load(efs, "de", "assets/i18n/de.yaml")

title := translations.T("de", "page.home.title")
tplData := translations.TMap("en", "nav.home", "nav.about", "action.save")
```

---

### [`logging`](logging/) — Structured Logging

JSON-based structured logging with configurable log levels and HTTP middleware.

| Function | Description |
|----------|-------------|
| `NewJsonLogger()` | Creates a JSON logger (level via `LOGGING_LEVEL` env var) |
| `WithLogging(logger, handler)` | HTTP middleware that logs requests with duration |

```go
import "github.com/andygeiss/cloud-native-utils/logging"

logger := logging.NewJsonLogger() // Reads LOGGING_LEVEL from env
logger.Info("server started", "port", 8080)

// HTTP middleware
http.Handle("/api", logging.WithLogging(logger, handler))
```

---

### [`messaging`](messaging/) — Pub/Sub Messaging

Decouples services using publish/subscribe patterns with Kafka or in-memory dispatchers.

| Type | Description |
|------|-------------|
| `Dispatcher` | Interface with `Publish` and `Subscribe` methods |
| `NewExternalDispatcher()` | Kafka-backed dispatcher (uses `KAFKA_BROKERS` env var) |
| `NewInternalDispatcher()` | In-memory dispatcher for local/test scenarios |
| `Message` | Message struct with topic, data, and state |

```go
import "github.com/andygeiss/cloud-native-utils/messaging"

// Production: Kafka
dispatcher := messaging.NewExternalDispatcher()

// Development/Testing: In-memory
dispatcher := messaging.NewInternalDispatcher()

// Subscribe
dispatcher.Subscribe(ctx, "orders", func(ctx context.Context, msg messaging.Message) (messaging.MessageState, error) {
    return messaging.MessageStateCompleted, nil
})

// Publish
dispatcher.Publish(ctx, messaging.NewMessage("orders", orderData))
```

---

### [`redirecting`](redirecting/) — PRG Pattern & HTMX Support

Implements Post/Redirect/Get pattern with automatic HTMX compatibility.

| Function | Description |
|----------|-------------|
| `Redirect(w, r, target)` | HTMX-aware redirect (HX-Redirect or 303 See Other) |
| `RedirectWithMessage(w, r, target, key, value)` | Redirect with query parameter |
| `WithPRG(handler)` | Middleware that auto-converts redirects for HTMX |

```go
import "github.com/andygeiss/cloud-native-utils/redirecting"

func CreateHandler(w http.ResponseWriter, r *http.Request) {
    redirecting.Redirect(w, r, "/items") // Works with both HTMX and standard requests
}

// Or use middleware
mux.Handle("/", redirecting.WithPRG(handler))
```

---

### [`resource`](resource/) — Generic Data Access

Unified CRUD interface with multiple storage backends.

| Type | Description |
|------|-------------|
| `Access[K, V]` | Generic interface: `Create`, `Read`, `ReadAll`, `Update`, `Delete` |
| `NewInMemoryAccess[K, V]()` | Thread-safe in-memory implementation |
| `NewSqliteAccess[K, V](db)` | SQLite-backed implementation |
| `NewJsonFileAccess[K, V](path)` | JSON file persistence |
| `NewYamlFileAccess[K, V](path)` | YAML file persistence |
| `NewMockAccess[K, V]()` | Test mock with builder pattern |

```go
import "github.com/andygeiss/cloud-native-utils/resource"

// In-memory storage
store := resource.NewInMemoryAccess[string, Task]()

// CRUD operations
store.Create(ctx, "task-1", task)
task, err := store.Read(ctx, "task-1")
tasks, err := store.ReadAll(ctx)
store.Update(ctx, "task-1", updatedTask)
store.Delete(ctx, "task-1")

// Mock for testing
mock := resource.NewMockAccess[string, Task]().
    WithReadFn(func(ctx context.Context, key string) (*Task, error) {
        return &Task{ID: key}, nil
    })
```

---

### [`security`](security/) — Cryptography & HTTP Security

Comprehensive security utilities for encryption, authentication, and secure HTTP handling.

| Category | Functions |
|----------|-----------|
| **Encryption** | `Encrypt(plaintext, key)`, `Decrypt(ciphertext, key)` — AES-256-GCM |
| **Key Generation** | `GenerateKey()` — 256-bit key, `GenerateID()` — hex-encoded ID |
| **PKCE** | `GeneratePKCE()` — OAuth 2.0 code verifier and challenge |
| **Hashing** | `Hash(tag, data)` — HMAC-SHA512/256 |
| **Password** | `Password(plaintext)`, `IsPasswordValid(hash, plaintext)` — bcrypt |
| **Env Parsing** | `ParseStringOrDefault`, `ParseIntOrDefault`, `ParseBoolOrDefault`, `ParseDurationOrDefault`, `ParseFloatOrDefault` |
| **HTTP Server** | `NewServer(mux)` — Configured with timeouts from env vars |
| **HTTP Client** | `NewClient()`, `NewClientWithTLS(cert, key, ca)` — Secure clients |
| **ServeMux** | `NewServeMux(ctx, efs)` — Pre-configured with `/liveness`, `/readiness`, `/static/`, and OIDC endpoints |
| **Sessions** | `ServerSessions` — Thread-safe session management |
| **Auth Middleware** | `WithAuth(sessions, handler)` — Extracts session/claims into context |
| **OIDC** | `IdentityProvider` — Complete OpenID Connect login/logout/callback flow |

```go
import "github.com/andygeiss/cloud-native-utils/security"

// Encryption
key := security.GenerateKey()
ciphertext := security.Encrypt([]byte("secret"), key)
plaintext, _ := security.Decrypt(ciphertext, key)

// Configuration
port := security.ParseIntOrDefault("PORT", 8080)
timeout := security.ParseDurationOrDefault("TIMEOUT", 5*time.Second)

// HTTP Server with health probes and OIDC
mux, sessions := security.NewServeMux(ctx, efs)
server := security.NewServer(mux)
server.ListenAndServe()
```

---

### [`service`](service/) — Lifecycle & Context Management

Cloud-native service orchestration with signal handling and context-aware functions.

| Type/Function | Description |
|---------------|-------------|
| `Function[IN, OUT]` | Generic function type with context support |
| `Context()` | Creates context that listens for OS signals (SIGTERM, SIGINT, etc.) |
| `RegisterOnContextDone(ctx, fn)` | Registers cleanup function for graceful shutdown |
| `Wrap[IN, OUT](fn)` | Converts simple function to context-aware function |

```go
import "github.com/andygeiss/cloud-native-utils/service"

func main() {
    ctx, cancel := service.Context()
    defer cancel()

    server := startServer()

    // Register graceful shutdown
    service.RegisterOnContextDone(ctx, func() {
        server.Shutdown(context.Background())
    })

    server.ListenAndServe()
}
```

---

### [`stability`](stability/) — Resilience Patterns

Production-ready stability patterns for fault-tolerant services.

| Function | Description |
|----------|-------------|
| `Breaker[IN, OUT](fn, threshold)` | Circuit breaker with exponential backoff |
| `Retry[IN, OUT](fn, maxRetries, delay)` | Automatic retry with configurable delay |
| `Throttle[IN, OUT](fn, maxTokens, refill, duration)` | Token bucket rate limiting |
| `Debounce[IN, OUT](fn, duration)` | Coalesces rapid calls into single execution |
| `Timeout[IN, OUT](fn, duration)` | Enforces execution time limits |

```go
import "github.com/andygeiss/cloud-native-utils/stability"

// Wrap function with multiple stability patterns
fn := fetchData
fn = stability.Retry(fn, 3, time.Second)           // Retry up to 3 times
fn = stability.Timeout(fn, 5*time.Second)          // 5 second timeout
fn = stability.Breaker(fn, 5)                      // Open circuit after 5 failures
fn = stability.Throttle(fn, 100, 10, time.Second)  // Rate limit: 100 tokens, refill 10/sec

result, err := fn(ctx, input)
```

---

### [`templating`](templating/) — Template Engine

Simple wrapper around Go's `text/template` with embedded filesystem support.

| Type/Method | Description |
|-------------|-------------|
| `Engine` | Template engine with embedded FS support |
| `NewEngine(efs)` | Creates engine with embedded filesystem |
| `Parse(patterns...)` | Parses templates using glob patterns |
| `Render(w, name, data)` | Executes named template with data |

```go
import "github.com/andygeiss/cloud-native-utils/templating"

//go:embed assets/templates/*
var efs embed.FS

engine := templating.NewEngine(efs)
engine.Parse("assets/templates/*.tmpl", "assets/templates/**/*.tmpl")

engine.Render(w, "index.tmpl", map[string]any{
    "Title": "Home",
    "Items": items,
})
```

---

## **Getting Started**

The repository [cloud-native-app](https://github.com/andygeiss/cloud-native-app)
provides a complete application scaffold that demonstrates how these modules
work together in a production-ready service.

```bash
# Generate a new cloud-native application
go install github.com/andygeiss/cloud-native-app@latest
cloud-native-app myservice
```

---

## **Dependencies**

| Dependency | Purpose |
|------------|---------|
| `github.com/coreos/go-oidc/v3` | OpenID Connect client |
| `github.com/segmentio/kafka-go` | Kafka messaging |
| `golang.org/x/crypto` | bcrypt password hashing |
| `golang.org/x/oauth2` | OAuth2 client |
| `gopkg.in/yaml.v3` | YAML parsing for i18n |
| `modernc.org/sqlite` | Pure-Go SQLite driver |

---

## **License**

MIT License — see [LICENSE](LICENSE) for details.
