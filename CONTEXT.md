# CONTEXT.md - Cloud Native Utils Codebase Guide

## Quick Reference (Agent Scan Block)

```
PROJECT: cloud-native-utils
MODULE: github.com/andygeiss/cloud-native-utils
LANGUAGE: Go 1.25.4+
LICENSE: MIT

CORE_TYPE: type Function[IN, OUT any] func(ctx context.Context, in IN) (out OUT, err error)

PACKAGES: assert, consistency, efficiency, env, event, extensibility, logging, mcp, messaging, resource, security, service, slices, stability, templating, web

TOOLING:
- LINTER: golangci-lint (config: .golangci.yml)
- TASK_RUNNER: just (config: .justfile)
- COMMANDS: just test, just lint, just benchmark

PATTERNS: Decorator, Strategy, Factory, Observer/PubSub, Adapter

RULES:
- CONTEXT_FIRST: Always pass context.Context as first parameter
- ERROR_LAST: Return error as last value in tuple
- ERROR_CONSTANTS: Use const/var for expected errors
- CONCURRENCY: sync.RWMutex for shared state, buffered channels for async
- CONSTRUCTOR: NewTypeName() for factories
- TEST_NAMING: Test_<Function>_With_<Condition>_Should_<Behavior>
- RECEIVER: Single lowercase letter (a *Type)
```

---

## 1. Project Overview

### 1.1 Purpose

A modular Go library providing reusable utilities for building cloud-native applications. Each package is independently importable with no monolithic framework.

### 1.2 Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/coreos/go-oidc/v3` | v3.17.0 | OIDC authentication |
| `github.com/segmentio/kafka-go` | v0.4.49 | Kafka messaging |
| `golang.org/x/crypto` | v0.46.0 | Cryptographic operations |
| `golang.org/x/oauth2` | v0.34.0 | OAuth2 flows |
| `gopkg.in/yaml.v3` | v3.0.1 | YAML parsing |
| `modernc.org/sqlite` | v1.40.1 | SQLite database |

### 1.3 Package Map

| Package | Domain | Key Types/Functions |
|---------|--------|---------------------|
| `assert` | Testing | `That()` |
| `consistency` | Event Sourcing | `Logger[K,V]`, `Event[K,V]`, `JsonFileLogger` |
| `efficiency` | Performance | `Generate`, `Merge`, `Split`, `Process`, `WithCompression` |
| `env` | Configuration | `Get[T]` |
| `event` | Domain Events | `Event`, `EventPublisher`, `EventSubscriber` |
| `extensibility` | Plugins | `LoadPlugin()` |
| `logging` | Observability | `NewJsonLogger()` |
| `mcp` | AI Integration | `Server`, `Tool`, `ToolHandler`, `NewServer()` |
| `messaging` | Pub/Sub | `Dispatcher`, `Message`, `NewMessage()` |
| `resource` | Persistence | `Access[K,V]`, `IndexedAccess[K,V]` |
| `security` | Crypto | `Encrypt`, `Decrypt`, `Hash`, `Password`, `GenerateKey` |
| `service` | Lifecycle | `Function[IN,OUT]`, `Context()`, `Wrap()` |
| `slices` | Collections | `Map`, `Filter`, `Unique`, `Contains` |
| `stability` | Resilience | `Retry`, `Breaker`, `Throttle`, `Debounce` |
| `templating` | Views | `Engine` |
| `web` | HTTP | `NewServer`, `NewClient`, `NewServeMux`, `WithAuth`, `WithBearerAuth` |

---

## 2. Core Type System

### 2.1 The Function Type (Foundation)

All service functions use this generic type signature:

```go
// FILE: service/function.go
// Function gathers together things that change for the same reasons.
// A context must be handled to be "cloud native" because it allows
// propagation of deadlines, cancellation signals, and other request-scoped values
// across API boundaries and between processes.
type Function[IN, OUT any] func(ctx context.Context, in IN) (out OUT, err error)
```

**Key Properties:**
- Generic with `[IN, OUT any]` type parameters
- Context-first for cancellation/timeout propagation
- Returns tuple of (result, error)

**Usage Example:**
```go
// Define a typed service function
var getUserByID service.Function[string, *User] = func(ctx context.Context, id string) (*User, error) {
    // implementation
    return user, nil
}
```

### 2.2 Generic Interface Pattern

Interfaces use type parameters for type-safe operations:

```go
// FILE: resource/access.go
// Access specifies the CRUD operations for a resource using generics.
// It supports context.Context for cancellation and timeouts.
type Access[K, V any] interface {
    Create(ctx context.Context, key K, value V) error
    Read(ctx context.Context, key K) (*V, error)
    ReadAll(ctx context.Context) ([]V, error)
    Update(ctx context.Context, key K, value V) error
    Delete(ctx context.Context, key K) error
}
```

### 2.3 Type Constraints

| Constraint | When to Use | Example |
|------------|-------------|---------|
| `any` | Maximum flexibility | `func Generate[T any](in ...T)` |
| `comparable` | Map keys or equality | `type IndexedAccess[K comparable, V any]` |

---

## 3. Architectural Patterns

### 3.1 Decorator/Wrapper Pattern (Stability)

Wrap functions to add cross-cutting concerns while preserving the `Function[IN,OUT]` signature.

```go
// FILE: stability/retry.go
// Retry wraps a given function (`fn`) to retry its execution upon failure.
// The function will be retried up to `maxRetries` times with a delay of `delay` between retries.
// If the context is canceled during retries, it stops immediately and returns the context error.
func Retry[IN, OUT any](fn service.Function[IN, OUT], maxRetries int, delay time.Duration) service.Function[IN, OUT] {
    return func(ctx context.Context, in IN) (out OUT, err error) {
        if ctx.Err() != nil {
            return out, ctx.Err()
        }
        for retries := 0; ; retries++ {
            res, err := fn(ctx, in)
            if err == nil || retries >= maxRetries {
                return res, err
            }
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return out, ctx.Err()
            }
        }
    }
}
```

**Pattern Template:**
```go
func Decorator[IN, OUT any](fn service.Function[IN, OUT], /* config params */) service.Function[IN, OUT] {
    // 1. Initialize decorator state (closured variables)
    var mutex sync.Mutex

    return func(ctx context.Context, in IN) (out OUT, err error) {
        // 2. Check context first
        if ctx.Err() != nil {
            return out, ctx.Err()
        }

        // 3. Pre-processing logic

        // 4. Call wrapped function
        res, err := fn(ctx, in)

        // 5. Post-processing logic

        return res, err
    }
}
```

**Available Decorators:**

| Function | Purpose | Parameters |
|----------|---------|------------|
| `Retry` | Retry on failure | `maxRetries int, delay time.Duration` |
| `Breaker` | Circuit breaker with exponential backoff | `threshold int` |
| `Throttle` | Token bucket rate limiting | `maxTokens, refill uint, duration time.Duration` |
| `Debounce` | Collapse rapid calls | `duration time.Duration` |

### 3.2 Strategy Pattern (Resource Access)

Multiple implementations of the `Access[K,V]` interface:

```go
// In-memory implementation
store := resource.NewInMemoryAccess[string, User]()

// JSON file implementation
store := resource.NewJsonFileAccess[string, User]("users.json")

// YAML file implementation
store := resource.NewYamlFileAccess[string, User]("users.yaml")

// SQLite implementation
store := resource.NewSqliteAccess[string, User](db, "users")

// PostgreSQL implementation
store := resource.NewPostgresAccess[string, User](db)
_ = store.Init(ctx) // Creates kv_store table and index
```

### 3.3 Adapter Pattern (IndexedAccess)

Wrap existing Access to add secondary indexing:

```go
// FILE: resource/indexed_access.go
// IndexedAccess wraps a resource.Access and maintains secondary indexes.
// It supports both unique and non-unique indexes (stored as lists of keys).
type IndexedAccess[K comparable, V any] struct {
    access     Access[K, V]
    indexes    map[string]map[string][]K
    indexFuncs map[string]IndexFunc[V]
    mu         sync.RWMutex
}

func NewIndexedAccess[K comparable, V any](access Access[K, V]) *IndexedAccess[K, V] {
    return &IndexedAccess[K, V]{
        access:     access,
        indexes:    make(map[string]map[string][]K),
        indexFuncs: make(map[string]IndexFunc[V]),
    }
}
```

### 3.4 Observer/Pub-Sub Pattern (Messaging)

```go
// FILE: messaging/dispatcher.go
// Dispatcher is an interface for a message dispatcher.
type Dispatcher interface {
    Publish(ctx context.Context, message Message) error
    Subscribe(ctx context.Context, topic string, fn service.Function[Message, MessageState]) error
}

// Message is a struct that represents a message.
type Message struct {
    Data  []byte       `json:"data"`
    State MessageState `json:"state"`
    Topic string       `json:"topic"`
}

// NewMessage creates a new message.
func NewMessage(topic string, data []byte) Message {
    return Message{
        Data:  data,
        State: MessageStateCreated,
        Topic: topic,
    }
}
```

**Implementations:**
- `NewInternalDispatcher()` - in-memory, goroutine-based
- `NewExternalDispatcher()` - Kafka-backed

### 3.5 Factory Pattern

All complex types use `New*` constructors:

```go
// Server factory (web package)
server := web.NewServer(mux)

// Client factory (web package)
client := web.NewClient()
client := web.NewClientWithTLS(certFile, keyFile, caFile)

// Message factory
msg := messaging.NewMessage(topic, data)

// Session store factory (web package)
sessions := web.NewServerSessions()

// Logger factory
logger := logging.NewJsonLogger()
```

---

## 4. Naming Conventions

### 4.1 Package Names

| Rule | Example | Anti-pattern |
|------|---------|--------------|
| Lowercase | `stability` | `Stability` |
| Single-word | `messaging` | `message_broker` |
| Domain-oriented | `security` | `utils` |

### 4.2 File Names

| Pattern | Example |
|---------|---------|
| Feature file | `retry.go`, `encrypt.go` |
| Test file | `retry_test.go` |
| Multi-word | `json_file_access.go`, `server_sessions.go` |

### 4.3 Type Names

| Category | Pattern | Examples |
|----------|---------|----------|
| Exported struct | PascalCase | `ServerSession`, `IndexedAccess` |
| Generic type | `Name[K, V any]` | `Function[IN, OUT any]`, `Access[K, V any]` |
| Interface | Noun | `Dispatcher`, `Logger`, `Access` |
| Unexported struct | camelCase | `gzipResponseWriter` |

### 4.4 Function Names

| Category | Pattern | Examples |
|----------|---------|----------|
| Constructor | `New<Type>` | `NewServer()`, `NewClient()`, `NewMessage()` |
| Predicate | `Is<Condition>` | `IsPasswordValid()` |
| Action | Verb phrase | `Encrypt()`, `Decrypt()`, `Parse()` |
| Decorator | Verb | `Retry()`, `Breaker()`, `Throttle()` |

### 4.5 Variable Names

| Category | Pattern | Examples |
|----------|---------|----------|
| Receiver | Single letter `a` | `func (a *Engine)`, `func (a *IndexedAccess)` |
| Context | `ctx` | `func(ctx context.Context, ...)` |
| Error | `err` | `res, err := fn(ctx, in)` |
| Mutex | `mutex` or `mu` | `var mutex sync.RWMutex` |
| Channel | `<purpose>Ch` | `errCh`, `stateCh`, `eventCh` |

### 4.6 Constants

**Error Constants Pattern:**
```go
// FILE: resource/access.go
const (
    ErrorResourceAlreadyExists = "Resource already exists"
    ErrorResourceNotFound      = "Resource not found"
)

// FILE: stability/breaker.go
var ErrorBreakerServiceUnavailable = errors.New("Service unavailable")

// FILE: stability/throttle.go
var ErrorThrottleTooManyCalls = errors.New("Too many calls")
```

**Enum Pattern:**
```go
// FILE: messaging/dispatcher.go
type MessageState int

const (
    MessageStateCreated MessageState = iota
    MessageStateCompleted
    MessageStateFailed
)
```

### 4.7 Test Function Names

**Pattern:** `Test_<Function>_With_<Condition>_Should_<Behavior>`

```go
// FILE: stability/retry_test.go
func Test_Retry_With_AlwaysFailingFunction_Should_ReturnError(t *testing.T)
func Test_Retry_With_SuccessAfterRetries_Should_ReturnResult(t *testing.T)
func Test_Retry_With_SuccessfulFunction_Should_ReturnResult(t *testing.T)
```

---

## 5. Coding Conventions

### 5.1 Context Handling

**Rule:** Context is ALWAYS the first parameter.

```go
// CORRECT
func Create(ctx context.Context, key K, value V) error

// Check context before expensive operations
if ctx.Err() != nil {
    return out, ctx.Err()
}
```

**Signal-aware context creation:**
```go
// FILE: service/context.go
// Context creates a new context with a cancel function that listens for
// SIGTERM, SIGINT, SIGQUIT, and SIGKILL signals.
func Context() (ctx context.Context, cancel context.CancelFunc) {
    return signal.NotifyContext(
        context.Background(),
        syscall.SIGTERM,  // Kubernetes graceful stop
        syscall.SIGINT,   // Terminal interrupt
        syscall.SIGQUIT,  // Core dump
        syscall.SIGKILL,  // Force kill
    )
}
```

### 5.2 Error Handling

**Rules:**
1. Return error as last value
2. Use error constants for expected failures
3. Check errors immediately

```go
// Error constants pattern
const (
    ErrorResourceNotFound = "Resource not found"
)

// Error variables for errors.Is() support
var ErrorBreakerServiceUnavailable = errors.New("Service unavailable")

// Checking pattern
res, err := fn(ctx, in)
if err != nil {
    return out, err
}
```

### 5.3 Concurrency Patterns

**Mutex for shared state:**
```go
type IndexedAccess[K comparable, V any] struct {
    mu sync.RWMutex  // Use RWMutex for read-heavy workloads
}

// Read operation - use RLock
func (a *IndexedAccess[K, V]) FindByIndex(ctx context.Context, indexName, indexKey string) ([]V, error) {
    a.mu.RLock()
    idx, ok := a.indexes[indexName]
    a.mu.RUnlock()
    // ...
}

// Write operation - use Lock
func (a *IndexedAccess[K, V]) Create(ctx context.Context, key K, value V) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    // ...
}
```

**Buffered channels for async results:**
```go
ch := make(chan response, 1)  // Buffered to prevent goroutine leak
```

**Select with context:**
```go
select {
case <-time.After(delay):
    // timeout action
case <-ctx.Done():
    return out, ctx.Err()
case res := <-ch:
    return res.result, res.err
}
```

### 5.4 Environment Configuration

**Pattern:** Use the `env` package for type-safe environment variable parsing with defaults.

```go
// FILE: env/env.go
// Get retrieves an environment variable and parses it to the type of the default value.
// Supported types: bool, int, float64, string, time.Duration
func Get[T any](key string, def T) T

// Usage examples
timeout := env.Get("SERVER_TIMEOUT", 5*time.Second)
maxRetries := env.Get("MAX_RETRIES", 3)
debug := env.Get("DEBUG", false)
rate := env.Get("RATE_LIMIT", 1.5)
name := env.Get("APP_NAME", "default")
```

**Server configuration example:**
```go
// FILE: web/server.go
func NewServer(mux *http.ServeMux) *http.Server {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    return &http.Server{
        Addr:              ":" + port,
        Handler:           mux,
        IdleTimeout:       env.Get("SERVER_IDLE_TIMEOUT", 5*time.Second),
        MaxHeaderBytes:    1 << 20,
        ReadHeaderTimeout: env.Get("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
        ReadTimeout:       env.Get("SERVER_READ_TIMEOUT", 5*time.Second),
        WriteTimeout:      env.Get("SERVER_WRITE_TIMEOUT", 5*time.Second),
    }
}
```

**Environment Variables:**

| Variable | Default | Package |
|----------|---------|---------|
| `PORT` | `8080` | web |
| `LOGGING_LEVEL` | `INFO` | logging |
| `CLIENT_TIMEOUT` | `5s` | web |
| `SERVER_*_TIMEOUT` | `5s` | web |
| `KAFKA_BROKERS` | - | messaging |

### 5.5 Package Documentation

```go
// FILE: <package>/<package>.go
// Package <name> <brief description of what it provides>.
package <name>
```

Example:
```go
// Package stability ensures service robustness with mechanisms like circuit breakers,
// retries, throttling, and debouncing for resilient distributed systems.
package stability
```

### 5.6 Import Organization

**Standard order:**
1. Standard library
2. External dependencies
3. Internal packages

```go
import (
    // Standard library
    "context"
    "sync"
    "time"

    // External dependencies
    "github.com/segmentio/kafka-go"

    // Internal packages
    "github.com/andygeiss/cloud-native-utils/service"
)
```

---

## 6. Security Patterns

### 6.1 Encryption (AES-GCM)

```go
// FILE: security/encrypt.go
// Encrypt takes an input byte slice (plaintext) and encrypts it using AES-GCM.
func Encrypt(plaintext []byte, key [32]byte) (ciphertext []byte) {
    block, _ := aes.NewCipher(key[:])
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    _, _ = io.ReadFull(rand.Reader, nonce)
    return gcm.Seal(nonce, nonce, plaintext, nil)
}
```

### 6.2 Secure Cookie Settings

```go
cookie := http.Cookie{
    Name:     "sid",
    Value:    sessionID,
    Path:     "/",
    HttpOnly: true,                    // Prevents JavaScript access
    Secure:   true,                    // HTTPS only
    SameSite: http.SameSiteLaxMode,    // CSRF protection
}
```

### 6.3 Security Headers

```go
w.Header().Set("Cache-Control", "no-store")
w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
```

---

## 7. Testing Patterns

### 7.1 Test File Organization

- Tests co-located: `<file>.go` + `<file>_test.go`
- Black-box testing: `package stability_test`
- Mocks in separate file: `mocks_test.go`

### 7.2 Assert Helper

```go
// FILE: assert/that.go
// That is a utility function to assert that two values are equal during a test.
func That(t *testing.T, desc string, got, expected any) {
    if !reflect.DeepEqual(got, expected) {
        t.Errorf("%s, but got %v (expected %v)", desc, got, expected)
    }
}

// Usage
assert.That(t, "err must be correct", err.Error(), "error")
assert.That(t, "result must be correct", res, 42)
assert.That(t, "err must be nil", err == nil, true)
```

### 7.3 AAA Pattern (Arrange-Act-Assert)

```go
// FILE: stability/retry_test.go
func Test_Retry_With_SuccessAfterRetries_Should_ReturnResult(t *testing.T) {
    // Arrange
    fn := stability.Retry(mockFailsTimes(2), 3, 10*time.Millisecond)

    // Act
    res, err := fn(context.Background(), 42)

    // Assert
    assert.That(t, "err must be nil", err == nil, true)
    assert.That(t, "result must be correct", res, 42)
}
```

### 7.4 Mock Functions Pattern

```go
// FILE: stability/mocks_test.go
func mockAlwaysFails() service.Function[int, int] {
    return func() service.Function[int, int] {
        return func(ctx context.Context, in int) (out int, err error) {
            return out, errors.New("error")
        }
    }()
}

func mockAlwaysSucceeds() service.Function[int, int] {
    return func() service.Function[int, int] {
        return func(ctx context.Context, in int) (int, error) {
            return 42, nil
        }
    }()
}

func mockSlow(duration time.Duration) service.Function[int, int] {
    return func(ctx context.Context, in int) (int, error) {
        select {
        case <-ctx.Done():
            return 0, ctx.Err()
        case <-time.After(duration):
            return in * 2, nil
        }
    }
}

func mockFailsTimes(n int) service.Function[int, int] {
    return func() service.Function[int, int] {
        var count int
        var mutex sync.Mutex
        return func(ctx context.Context, in int) (out int, err error) {
            mutex.Lock()
            defer mutex.Unlock()
            if count >= n {
                return 42, nil
            }
            count++
            return out, errors.New("error")
        }
    }()
}
```

---

## 8. HTTP Middleware Pattern

### 8.1 Standard Middleware Signature

```go
func WithMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Pre-processing

        next.ServeHTTP(w, r)

        // Post-processing
    })
}
```

### 8.2 ResponseWriter Wrapper Pattern

```go
// FILE: efficiency/middleware.go
type gzipResponseWriter struct {
    http.ResponseWriter
    gzw         *gzip.Writer
    wroteHeader bool
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
    return &gzipResponseWriter{
        ResponseWriter: w,
        gzw:            gzip.NewWriter(w),
    }
}

func (a *gzipResponseWriter) Write(p []byte) (int, error) {
    if !a.wroteHeader {
        a.WriteHeader(http.StatusOK)
    }
    return a.gzw.Write(p)
}

func (a *gzipResponseWriter) WriteHeader(code int) {
    if a.wroteHeader {
        return
    }
    h := a.Header()
    h.Set("Content-Encoding", "gzip")
    h.Del("Content-Length")
    h.Add("Vary", "Accept-Encoding")
    a.ResponseWriter.WriteHeader(code)
    a.wroteHeader = true
}
```

### 8.3 Compression Middleware Example

```go
// FILE: efficiency/middleware.go
func WithCompression(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if encHeader := r.Header.Get("Accept-Encoding"); !strings.Contains(encHeader, "gzip") ||
            r.Header.Get("Range") != "" ||
            r.Method == http.MethodHead {
            next.ServeHTTP(w, r)
            return
        }

        gzw := newGzipResponseWriter(w)
        defer gzw.Close()

        next.ServeHTTP(gzw, r)
    })
}
```

---

## 9. Channel Pipeline Pattern

### 9.1 Generator Pattern

```go
// FILE: efficiency/generate.go
// Generate takes a variadic input of any type T and returns a read-only channel of type T.
func Generate[T any](in ...T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for _, val := range in {
            out <- val
        }
    }()
    return out
}
```

### 9.2 Fan-Out Worker Pool

```go
// FILE: efficiency/process.go
// Process concurrently processes items from the input channel using the provided function `fn`.
// It spawns a number of worker goroutines equal to the number of available CPU cores.
func Process[IN, OUT any](in <-chan IN, fn service.Function[IN, OUT]) (<-chan OUT, <-chan error) {
    out := make(chan OUT)
    errCh := make(chan error)
    ctx := context.Background()

    num := runtime.NumCPU()
    var wg sync.WaitGroup
    wg.Add(num)
    for range num {
        go func() {
            defer wg.Done()
            for val := range in {
                res, err := fn(ctx, val)
                if err != nil {
                    errCh <- err
                    return
                }
                out <- res
            }
        }()
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out, errCh
}
```

---

## 10. Code Generation Templates

### 10.1 New Package Template

```go
// FILE: <package>/<package>.go
// Package <package> provides <brief description>.
package <package>
```

### 10.2 New Service Function Template

```go
package <package>

import (
    "context"

    "github.com/andygeiss/cloud-native-utils/service"
)

// <FunctionName> <brief description>.
func <FunctionName>[IN, OUT any](/* config params */) service.Function[IN, OUT] {
    return func(ctx context.Context, in IN) (out OUT, err error) {
        if ctx.Err() != nil {
            return out, ctx.Err()
        }
        // Implementation
        return out, nil
    }
}
```

### 10.3 New Interface Implementation Template

```go
package <package>

import (
    "context"
    "sync"
)

type <implName> struct {
    mu sync.RWMutex
    // fields
}

// New<ImplName> creates a new <implName>.
func New<ImplName>() *<implName> {
    return &<implName>{
        // initialize fields
    }
}

func (a *<implName>) Method(ctx context.Context /* params */) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    // implementation
    return nil
}
```

### 10.4 New Test File Template

```go
package <package>_test

import (
    "context"
    "testing"

    "github.com/andygeiss/cloud-native-utils/assert"
    "github.com/andygeiss/cloud-native-utils/<package>"
)

func Test_<Function>_With_<Condition>_Should_<Behavior>(t *testing.T) {
    // Arrange

    // Act
    result, err := <package>.<Function>(context.Background(), input)

    // Assert
    assert.That(t, "err must be nil", err == nil, true)
    assert.That(t, "result must be correct", result, expected)
}
```

---

## 11. Common Pitfalls

### 11.1 Forgetting Context Check

```go
// BAD
func MyWrapper[IN, OUT any](fn service.Function[IN, OUT]) service.Function[IN, OUT] {
    return func(ctx context.Context, in IN) (out OUT, err error) {
        return fn(ctx, in)  // Missing context check!
    }
}

// GOOD
func MyWrapper[IN, OUT any](fn service.Function[IN, OUT]) service.Function[IN, OUT] {
    return func(ctx context.Context, in IN) (out OUT, err error) {
        if ctx.Err() != nil {
            return out, ctx.Err()
        }
        return fn(ctx, in)
    }
}
```

### 11.2 Incorrect Mutex Usage

```go
// BAD - Using Lock for reads
func (a *Store) Read(key string) *Value {
    a.mu.Lock()  // Should be RLock for reads!
    defer a.mu.Unlock()
    return a.data[key]
}

// GOOD
func (a *Store) Read(key string) *Value {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.data[key]
}
```

### 11.3 Leaking Goroutines

```go
// BAD - goroutine may leak if context canceled
go func() {
    for {
        // work forever
    }
}()

// GOOD - respect context cancellation
go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // work
        }
    }
}()
```

### 11.4 Unbuffered Channel Deadlock

```go
// BAD - may deadlock
ch := make(chan result)
go func() {
    res, err := fn(ctx, in)
    ch <- result{res, err}  // Blocks if no receiver
}()
// If we return early due to ctx.Done(), goroutine blocks forever

// GOOD - buffered channel
ch := make(chan result, 1)
```

### 11.5 Missing Channel Close

```go
// BAD - channel never closed, range will block forever
func Generate[T any](in ...T) <-chan T {
    out := make(chan T)
    go func() {
        for _, val := range in {
            out <- val
        }
        // Missing close(out)!
    }()
    return out
}

// GOOD
func Generate[T any](in ...T) <-chan T {
    out := make(chan T)
    go func() {
        defer close(out)
        for _, val := range in {
            out <- val
        }
    }()
    return out
}
```

---

## 12. Logging Pattern

```go
// FILE: logging/logger_json.go
// NewJsonLogger creates a new structured logger in JSON format.
func NewJsonLogger() *slog.Logger {
    var level slog.Leveler

    lvl := os.Getenv("LOGGING_LEVEL")
    lvl = strings.ToUpper(lvl)

    switch lvl {
    case "DEBUG":
        level = slog.LevelDebug
    case "ERROR":
        level = slog.LevelError
    case "INFO":
        level = slog.LevelInfo
    case "WARN":
        level = slog.LevelWarn
    default:
        level = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{Level: level}
    handler := slog.NewJSONHandler(os.Stdout, opts)

    return slog.New(handler)
}
```

**Usage:**
```go
logger := logging.NewJsonLogger()
logger.Info("http request handled", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
```

---

## 13. Key Files Reference

| Purpose | File Path |
|---------|-----------|
| Core Function type | `service/function.go` |
| Signal-aware context | `service/context.go` |
| Retry decorator | `stability/retry.go` |
| Circuit breaker | `stability/breaker.go` |
| Throttle decorator | `stability/throttle.go` |
| CRUD interface | `resource/access.go` |
| Indexed wrapper | `resource/indexed_access.go` |
| Pub/Sub interface | `messaging/dispatcher.go` |
| Channel pipeline | `efficiency/process.go`, `efficiency/generate.go` |
| HTTP middleware | `efficiency/middleware.go` |
| Environment config | `env/env.go` |
| Encryption | `security/encrypt.go` |
| Server factory | `web/server.go` |
| Auth middleware (session) | `web/middleware.go` |
| Auth middleware (bearer) | `web/middleware.go` |
| Identity provider | `web/identity_provider.go` |
| JSON logger | `logging/logger_json.go` |
| Test assertions | `assert/that.go` |
| Mock patterns | `stability/mocks_test.go` |
| MCP server | `mcp/server.go` |
| MCP tool types | `mcp/tool.go` |

---

## 14. MCP Server Pattern

### 14.1 Creating an MCP Server

```go
// FILE: mcp/server.go
server := mcp.NewServer("my-server", "1.0.0")

// Register tools
server.RegisterTool(tool)

// Start serving (STDIO transport)
ctx, cancel := service.Context()
defer cancel()
server.Serve(ctx)
```

### 14.2 Tool Definition

```go
// FILE: mcp/tool.go
// ToolHandler follows the Function[IN, OUT] pattern
type ToolHandler service.Function[ToolsCallParams, ToolsCallResult]

// Create a tool with schema and handler
schema := mcp.NewObjectSchema(
    map[string]mcp.Property{
        "name": mcp.NewStringProperty("The name to greet"),
    },
    []string{"name"},
)

handler := func(ctx context.Context, params mcp.ToolsCallParams) (mcp.ToolsCallResult, error) {
    name, _ := params.Arguments["name"].(string)
    return mcp.ToolsCallResult{
        Content: []mcp.ContentBlock{mcp.NewTextContent(fmt.Sprintf("Hello, %s!", name))},
    }, nil
}

tool := mcp.NewTool("greet", "Greets someone by name", schema, handler)
```

### 14.3 Schema Helpers

```go
// Object schema with properties
mcp.NewObjectSchema(properties map[string]Property, required []string) InputSchema

// Property types
mcp.NewStringProperty(description string) Property
mcp.NewNumberProperty(description string) Property
mcp.NewBooleanProperty(description string) Property
```

### 14.4 Response Helpers

```go
// Text content for tool results
mcp.NewTextContent(text string) ContentBlock

// JSON-RPC responses
mcp.NewResponse(id json.RawMessage, result any) Response
mcp.NewErrorResponse(id json.RawMessage, code int, message string) Response
```

### 14.5 Composing with Stability Patterns

```go
// Wrap tool handlers with stability patterns
handler := stability.Timeout(myHandler, 30*time.Second)
handler = stability.Retry(handler, 3, time.Second)

tool := mcp.NewTool("my-tool", "Description", schema, handler)
```

---

## 15. Web Package Pattern

### 15.1 Creating an HTTP Server

```go
// FILE: web/server.go
mux := http.NewServeMux()
server := web.NewServer(mux)
server.ListenAndServe()
```

### 15.2 Using NewServeMux with OIDC

```go
// FILE: web/mux.go
//go:embed assets
var efs embed.FS

ctx, cancel := service.Context()
defer cancel()

mux, sessions := web.NewServeMux(ctx, efs)
// mux now has /liveness, /readiness, /static/, and /auth/* endpoints
```

### 15.3 Session-Based Authentication Middleware

```go
// FILE: web/middleware.go
handler := web.WithAuth(sessions, func(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value(web.ContextEmail).(string)
    name := r.Context().Value(web.ContextName).(string)
    // Use authenticated user info
})
mux.HandleFunc("GET /protected", handler)
```

### 15.4 Bearer Token Authentication Middleware (MCP)

```go
// FILE: web/middleware.go
// Get OIDC verifier after identity provider is initialized via Login()
verifier := web.IdentityProvider.Verifier()

// Apply to MCP endpoint - returns JSON-RPC 2.0 errors on auth failure
handler := web.WithBearerAuth(verifier, func(w http.ResponseWriter, r *http.Request) {
    email := r.Context().Value(web.ContextEmail).(string)
    subject := r.Context().Value(web.ContextSubject).(string)
    // Use authenticated user info from Bearer token
})
mux.HandleFunc("POST /mcp", handler)
```

**Context values populated by WithBearerAuth:**
- `ContextEmail` - User's email from token claims
- `ContextName` - User's name from token claims
- `ContextSubject` - Subject (user ID) from token
- `ContextIssuer` - Token issuer URL
- `ContextVerified` - Email verification status

### 15.5 TLS Client

```go
// FILE: web/client.go
client := web.NewClient() // Basic client with timeout
client := web.NewClientWithTLS(certFile, keyFile, caFile) // mTLS client
```

---

## 16. Linting

### 16.1 Running the Linter

```bash
# Using just
just lint

# Directly
golangci-lint run ./...
```

### 16.2 Configuration

The `.golangci.yml` configuration:
- Uses `default: all` linters
- Disables overly strict linters (exhaustruct, varnamelen, mnd, etc.)
- Configures govet with all checks enabled
- Configures revive to allow missing package comments
- Configures gosec to allow 0755 directory permissions

### 16.3 Disabled Linters (Rationale)

| Linter | Reason |
|--------|--------|
| `exhaustruct` | Too verbose for tests |
| `varnamelen` | Too strict for idiomatic Go |
| `mnd` | Magic number detection too noisy |
| `wrapcheck` | Error wrapping too strict for simple cases |
| `forbidigo` | Print statements needed for CLI output |
| `err113` | Error comparison too strict |
