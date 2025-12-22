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

A collection of modular utilities for building cloud-native Go applications.
This repository is organized as small, focused Go packages (no monolithic framework), each designed to be used independently.

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
    - [Imaging](#imaging)
    - [Logging](#logging)
    - [Messaging](#messaging)
    - [Redirecting](#redirecting)
    - [Resource](#resource)
    - [Scheduling](#scheduling)
    - [Security](#security)
    - [Service](#service)
    - [Slices](#slices)
    - [Stability](#stability)
    - [Templating](#templating)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Assert** - Minimal test assertions (`assert.That`)
- **Consistency** - Transactional event log with JSON file persistence (`JsonFileLogger`)
- **Efficiency** - Channel helpers (`Generate`, `Merge`, `Split`, `Process`) and HTTP gzip (`WithCompression`)
- **Extensibility** - Dynamic plugin loading (`LoadPlugin`)
- **i18n** - Date/money formatting and YAML translations from `embed.FS` (`Translations`)
- **Imaging** - QR code generation (including Data URL output)
- **Logging** - Structured JSON logging via `log/slog`, plus an HTTP handler wrapper (`WithLogging`)
- **Messaging** - Publish/subscribe dispatcher (in-memory or Kafka-backed)
- **Redirecting** - PRG + HTMX-friendly redirects (`WithPRG`, `Redirect`, `RedirectWithMessage`)
- **Resource** - Generic CRUD access interface with multiple backends (memory/JSON/YAML/SQLite) and indexing
- **Scheduling** - Time/day primitives for booking systems (opening hours, slots, orphan gaps)
- **Security** - AES-GCM encryption, password hashing, env parsing helpers, OIDC identity provider helpers, and a configured HTTP server
- **Service** - Context helpers (signal-aware), function wrapper, and context-done hooks
- **Slices** - Generic slice helpers (`Map`, `Filter`, `Unique`, ...)
- **Stability** - Resilience wrappers for `service.Function` (breaker/retry/throttle/debounce/timeout)
- **Templating** - HTML templating engine on top of embedded filesystems

---

## Installation

```bash
go get github.com/andygeiss/cloud-native-utils
```

**Requirements:**
- Go 1.25.4 or later

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

Transactional log management with JSON file-based persistence.

```go
import "github.com/andygeiss/cloud-native-utils/consistency"

logger := consistency.NewJsonFileLogger[string, []byte]("./data/events.json")
defer logger.Close()

logger.WritePut("user:123", []byte("created"))
logger.WriteDelete("user:123")

events, errs := logger.ReadEvents()
for e := range events {
    _ = e // handle event
}
for err := range errs {
    if err != nil {
        // handle error
    }
}
```

### Efficiency

Utilities for concurrent stream processing.

```go
import (
    "context"
    "github.com/andygeiss/cloud-native-utils/efficiency"
)

// Generate a read-only channel from values
ch := efficiency.Generate(1, 2, 3, 4, 5)

// Merge multiple channels into one
ch1 := efficiency.Generate(1, 2, 3)
ch2 := efficiency.Generate(4, 5, 6)
merged := efficiency.Merge(ch1, ch2)
_ = merged

// Split a channel into multiple outputs (fan-out / work distribution)
input := efficiency.Generate(10, 11, 12, 13, 14)
workers := efficiency.Split(input, 3)
_ = workers

// Process items concurrently (worker count is based on NumCPU)
fn := func(ctx context.Context, in int) (int, error) { return in * 2, nil }
out, errCh := efficiency.Process(ch, fn)
_ = out
_ = errCh
```

Note: `Split` distributes items across outputs (it does not broadcast each item to every output).

### Extensibility

Dynamically load external Go plugins at runtime.

```go
import "github.com/andygeiss/cloud-native-utils/extensibility"

// Load a symbol from a plugin file
symbol, err := extensibility.LoadPlugin[func(string) string]("./plugins/myplugin.so", "MyFunction")
_ = symbol
_ = err
```

Note: Go plugins are platform-dependent (and not supported on all OS/architectures).

### i18n

YAML-based internationalization with embedded filesystem support.

```go
import "github.com/andygeiss/cloud-native-utils/i18n"

//go:embed translations/*.yaml
var translationsFS embed.FS

translations := i18n.NewTranslations()
_ = translations.Load(translationsFS, "en", "translations/en.yaml")
_ = translations.Load(translationsFS, "de", "translations/de.yaml")

// Get translated text
text := translations.T("en", "greeting.hello")
```

Date/money helpers are available as standalone functions too (e.g. `FormatDateISO`, `FormatMoney`).

### Imaging

QR code generation utilities.

```go
import "github.com/andygeiss/cloud-native-utils/imaging"

dataURL, err := imaging.GenerateQRCodeDataURL("https://example.com")
_ = dataURL
_ = err
```

### Logging

Structured JSON logging with HTTP middleware.

```go
import "github.com/andygeiss/cloud-native-utils/logging"

// Create a JSON logger (level configured via LOGGING_LEVEL)
logger := logging.NewJsonLogger()

// Wrap a handler func to emit structured request logs
handler := logging.WithLogging(logger, yourHandler)
```

### Messaging

Publish-subscribe patterns for decoupling services.

```go
import "github.com/andygeiss/cloud-native-utils/messaging"

// Create a dispatcher for internal messaging
dispatcher := messaging.NewInternalDispatcher()

// Subscribe to a topic
_ = dispatcher.Subscribe(ctx, "user.created", func(ctx context.Context, msg messaging.Message) (messaging.MessageState, error) {
    return messaging.MessageStateCreated, nil
})

// Publish a message
_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", payload))
```

For Kafka-backed messaging, use `messaging.NewExternalDispatcher()` and set `KAFKA_BROKERS`.

### Redirecting

HTMX-compatible HTTP redirects for POST/PUT/DELETE requests.

```go
import "github.com/andygeiss/cloud-native-utils/redirecting"

// Wrap a handler tree to translate redirect responses for HTMX (PRG support)
handler := redirecting.WithPRG(yourHandler)
```

For direct redirects, use `redirecting.Redirect` or `redirecting.RedirectWithMessage`.

### Resource

Generic CRUD interface with multiple backend implementations.

```go
import "github.com/andygeiss/cloud-native-utils/resource"

// In-memory storage
store := resource.NewInMemoryAccess[string, User]()

// JSON file storage
store := resource.NewJsonFileAccess[string, User]("users.json")

// SQLite storage
store := resource.NewSqliteAccess[string, User](db)

// CRUD operations
_ = store.Create(ctx, "user-1", user)
userPtr, err := store.Read(ctx, "user-1")
_ = userPtr
_ = store.Update(ctx, "user-1", updatedUser)
_ = store.Delete(ctx, "user-1")
users, err := store.ReadAll(ctx)
_ = users
_ = err
```

You can add a secondary index with `resource.NewIndexedAccess`.

### Scheduling

Time and scheduling primitives for booking systems.

```go
import "github.com/andygeiss/cloud-native-utils/scheduling"

open := scheduling.MustTimeOfDay(9, 0)
close := scheduling.MustTimeOfDay(17, 0)

day, err := scheduling.NewDayHours(scheduling.Monday, open, close)
_ = day
_ = err
```

### Security

Comprehensive security utilities for cloud-native applications.

```go
import (
    "net/http"
    "github.com/andygeiss/cloud-native-utils/security"
)

// AES-GCM encryption/decryption
key := security.GenerateKey()
ciphertext := security.Encrypt([]byte("secret"), key)
plaintext, err := security.Decrypt(ciphertext, key)
_ = plaintext
_ = err

// Password hashing + verification
hash, err := security.Password([]byte("p@ssw0rd"))
ok := security.IsPasswordValid(hash, []byte("p@ssw0rd"))
_ = ok
_ = err

// Generate secure IDs
id := security.GenerateID()

// PKCE for OAuth2
verifier, challenge := security.GeneratePKCE()

// Configured HTTP server (PORT and SERVER_*_TIMEOUT env vars)
mux := http.NewServeMux()
server := security.NewServer(mux)
_ = server
```

OIDC helpers are available via the `security.IdentityProvider` singleton (see `Login`, `Callback`, `Logout`).

### Service

Context-aware function wrappers with lifecycle management.

```go
import "github.com/andygeiss/cloud-native-utils/service"

ctx, cancel := service.Context()
defer cancel()

// Wrap a function with context support
fn := service.Wrap(func(in int) (int, error) {
    return in * 2, nil
})

// Register cleanup on context cancellation
service.RegisterOnContextDone(ctx, func() {
    // Cleanup logic
})
```

### Slices

Generic helpers for working with slices.

```go
import "github.com/andygeiss/cloud-native-utils/slices"

nums := []int{1, 2, 2, 3}
unique := slices.Unique(nums)
hasTwo := slices.Contains(nums, 2)
_ = unique
_ = hasTwo
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

- **Go** (1.25.4+) - Primary programming language
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
- Runs cleanly with `go test ./...`
- Has appropriate documentation

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
