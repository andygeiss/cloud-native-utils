# CONTEXT.md

## 1. Project purpose

Cloud Native Utils is a modular Go library providing reusable utilities for building cloud-native applications. Each package is designed to be independently importable and focused on a single responsibility—there is no monolithic framework.

The library addresses common cloud-native concerns: resilience patterns, structured logging, message dispatching, generic CRUD persistence, security primitives, and HTTP middleware. It serves as a utility toolkit that can be imported piecemeal into any Go project requiring these capabilities.

---

## 2. Technology stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.25.4+ |
| **Module path** | `github.com/andygeiss/cloud-native-utils` |
| **Build system** | Go modules (`go.mod`) |

### Direct dependencies

| Dependency | Purpose | Used in |
|------------|---------|---------|
| `github.com/coreos/go-oidc/v3` | OpenID Connect identity provider | `security` |
| `github.com/segmentio/kafka-go` | Kafka-backed messaging | `messaging` |
| `github.com/skip2/go-qrcode` | QR code generation | `imaging` |
| `golang.org/x/crypto` | bcrypt password hashing | `security` |
| `golang.org/x/oauth2` | OAuth2 client flows | `security` |
| `gopkg.in/yaml.v3` | YAML parsing | `resource` |
| `modernc.org/sqlite` | Pure-Go SQLite driver | `resource` |

---

## 3. High-level architecture

### Architectural style

**Modular library** – each top-level directory is an independent Go package with minimal cross-package dependencies. Consumers import only the packages they need.

### Core abstractions

| Abstraction | Package | Description |
|-------------|---------|-------------|
| `service.Function[IN, OUT]` | `service` | Context-aware function signature used across stability wrappers |
| `resource.Access[K, V]` | `resource` | Generic CRUD interface with multiple backends |
| `messaging.Dispatcher` | `messaging` | Publish-subscribe interface (internal & Kafka) |

### Package interaction patterns

```
┌──────────────┐
│   stability  │──wraps──▶ service.Function
└──────────────┘
       │
       ▼
┌──────────────┐
│   service    │ ◀── context, lifecycle, function types
└──────────────┘

┌──────────────┐     ┌──────────────┐
│   resource   │     │   messaging  │
│ Access[K,V]  │     │  Dispatcher  │
└──────────────┘     └──────────────┘
       │                    │
       ▼                    ▼
 InMemory / JSON /    Internal / External
 YAML / SQLite            (Kafka)
```

---

## 4. Directory structure (contract)

```
cloud-native-utils/
├── go.mod                 # Module definition
├── go.sum
├── README.md              # User-facing documentation
├── CONTEXT.md             # This file (AI/developer context)
├── LICENSE                # MIT
│
├── assert/                # Test assertion helpers
├── consistency/           # Transactional event log, JSON file logger
├── efficiency/            # Channel helpers, gzip middleware, sharding
├── extensibility/         # Dynamic plugin loading
├── imaging/               # QR code generation
├── logging/               # Structured JSON logging, HTTP middleware
├── messaging/             # Pub-sub dispatcher (internal, Kafka)
├── redirecting/           # HTMX-compatible PRG redirects
├── resource/              # Generic CRUD Access interface & backends
├── security/              # Encryption, hashing, OIDC, HTTP server
├── service/               # Context helpers, Function type, lifecycle
├── slices/                # Generic slice utilities
├── stability/             # Resilience patterns (breaker, retry, throttle, debounce, timeout)
└── templating/            # HTML template engine with fs.FS
```

### Rules for new code

| Code type | Location | Notes |
|-----------|----------|-------|
| New utility package | Top-level directory (e.g., `newpkg/`) | One package = one directory |
| Package doc comment | `<pkg>/<pkg>.go` (e.g., `logging/logging.go`) | Contains `// Package ...` |
| Tests | `<pkg>/<file>_test.go` | Same directory as source |
| Test fixtures | `<pkg>/testdata/` | Gitignored binaries, sample files |
| HTTP middleware | Existing package or new `<pkg>/middleware.go` | Follow `With*` naming |

---

## 5. Coding conventions

### 5.1 General

- **Small, focused packages** – each package addresses a single concern.
- **Generics where appropriate** – `Access[K, V]`, `Function[IN, OUT]`, slice helpers.
- **Context-first** – functions that may block accept `context.Context` as the first parameter.
- **No global state** – prefer constructor functions (`NewX`) returning struct pointers.
- **Filesystem-agnostic** – templates (templating.Engine) and static assets (security.NewServeMux) accept `fs.FS`, allowing `embed.FS`, `os.DirFS`, or custom implementations.

### 5.2 Naming

| Entity | Convention | Example |
|--------|------------|---------|
| Package | Lowercase, singular noun | `security`, `resource` |
| Exported function | PascalCase, verb phrase | `GenerateKey`, `WithLogging` |
| Interface | PascalCase noun | `Access`, `Dispatcher` |
| Constructor | `New<Type>` | `NewInMemoryAccess` |
| HTTP middleware | `With<Feature>` | `WithCompression`, `WithPRG` |
| Test file | `<source>_test.go` | `breaker_test.go` |
| Error sentinel | `Error<Description>` | `ErrorResourceNotFound` |
| Context key type | `ContextKey` (unexported string type) | `ContextSessionID` |

### 5.3 Error handling & logging

- **Sentinel errors** – define package-level error constants for expected failures (`ErrorResourceNotFound`).
- **Wrap with context** – use `fmt.Errorf("...: %w", err)` when adding context.
- **Check `ctx.Err()` early** – return immediately if context is cancelled.
- **Structured logging** – use `log/slog` via `logging.NewJsonLogger()`.
- **Log level via env** – `LOGGING_LEVEL` environment variable controls level.

### 5.4 Testing

- **Framework** – standard `testing` package only.
- **Assertions** – use `assert.That(t, desc, got, expected)` from this repo.
- **Table-driven tests** – preferred for multiple input/output scenarios.
- **Test file naming** – `<source>_test.go` in the same package.
- **Fixtures** – place in `testdata/` subdirectory.
- **No external test dependencies** – no testify, gomock, etc.

### 5.5 Formatting & linting

- **Formatter** – `gofmt` / `goimports` (standard Go tooling).
- **Linter** – Codacy integration; follow Go Report Card standards.
- **CI checks** – `go test ./...`, `go vet ./...`.

---

## 6. Cross-cutting concerns and reusable patterns

### Configuration

- Environment variables preferred (`os.Getenv`, `security.Getenv*` helpers).
- Common env vars: `PORT`, `LOGGING_LEVEL`, `KAFKA_BROKERS`, `SERVER_*_TIMEOUT`.

### HTTP middleware pattern

Middleware functions wrap `http.HandlerFunc`:

```go
func WithFeature(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // pre-processing
        next(w, r)
        // post-processing
    }
}
```

### Context propagation

- Use `service.Context()` for signal-aware root context.
- Pass `context.Context` as first parameter to blocking functions.
- Register cleanup via `service.RegisterOnContextDone(ctx, fn)`.

### Generic CRUD pattern

Implement `resource.Access[K, V]` interface:

```go
type Access[K, V any] interface {
    Create(ctx context.Context, key K, value V) error
    Read(ctx context.Context, key K) (*V, error)
    ReadAll(ctx context.Context) ([]V, error)
    Update(ctx context.Context, key K, value V) error
    Delete(ctx context.Context, key K) error
}
```

Provided implementations: `InMemoryAccess`, `JsonFileAccess`, `YamlFileAccess`, `SqliteAccess`, `IndexedAccess`, `MockAccess`.

### Resilience wrappers (stability package)

Wrap `service.Function[IN, OUT]` with resilience patterns:

| Pattern | Function | Purpose |
|---------|----------|---------|
| Circuit breaker | `stability.Breaker(fn, threshold)` | Stop calls after repeated failures |
| Retry | `stability.Retry(fn, attempts, delay)` | Retry transient failures |
| Throttle | `stability.Throttle(fn, maxConcurrent)` | Limit concurrent executions |
| Debounce | `stability.Debounce(fn, delay)` | Delay rapid successive calls |
| Timeout | `stability.Timeout(fn, duration)` | Enforce execution time limit |

### Messaging pattern

```go
dispatcher := messaging.NewInternalDispatcher() // or NewExternalDispatcher() for Kafka
dispatcher.Subscribe(ctx, topic, handler)
dispatcher.Publish(ctx, messaging.NewMessage(topic, payload))
```

### Approved vendor libraries

See `go.mod` for the authoritative list. Key patterns:

| Library | Usage pattern |
|---------|---------------|
| `gopkg.in/yaml.v3` | `yaml.Unmarshal` for config/translations |
| `modernc.org/sqlite` | Import as `_ "modernc.org/sqlite"` for driver registration |
| `golang.org/x/crypto/bcrypt` | `bcrypt.GenerateFromPassword`, `bcrypt.CompareHashAndPassword` |
| `github.com/coreos/go-oidc/v3` | OIDC provider verification |

---

## 7. Using this repo as a template

This repository is a **utility library**, not a project template. However, it establishes patterns that derived projects should follow:

### Invariants (preserve across projects)

- Context-first function signatures.
- Generic interfaces (`Access[K, V]`, `Function[IN, OUT]`).
- `With*` naming for HTTP middleware.
- Sentinel error constants.
- Standard `testing` package with `assert.That`.

### Customization points

- Add new packages at the top level for new concerns.
- Implement `resource.Access[K, V]` for new storage backends.
- Wrap `service.Function` with custom stability/observability logic.
- Extend `messaging.Dispatcher` for new message brokers.

### Recommended steps for new projects using this library

1. `go get github.com/andygeiss/cloud-native-utils`
2. Import only the packages you need (e.g., `security`, `stability`).
3. Follow the same conventions for your own code.
4. Use `service.Context()` for graceful shutdown.
5. Wrap handlers with provided middleware (`WithLogging`, `WithCompression`, `WithPRG`).

---

## 8. Key commands & workflows

| Task | Command |
|------|---------|
| Install dependencies | `go mod download` |
| Run all tests | `go test ./...` |
| Run tests with coverage | `go test -cover ./...` |
| Run tests verbose | `go test -v ./...` |
| Format code | `gofmt -w .` or `goimports -w .` |
| Vet code | `go vet ./...` |
| Build (library, no main) | `go build ./...` |
| Update dependencies | `go get -u ./... && go mod tidy` |

---

## 9. Important notes & constraints

### Security considerations

- **Secrets** – never log secrets; use `security.Getenv*` helpers.
- **Path sanitization** – use `logging.PathSanitizer` to mask sensitive URL segments.
- **OIDC** – configure via environment variables; see `security.IdentityProvider`.

### Platform assumptions

- **Go 1.25.4+** – uses generics and `log/slog`.
- **Plugin support** – `extensibility.LoadPlugin` only works on Linux/macOS (not Windows).
- **SQLite** – uses pure-Go `modernc.org/sqlite` (no CGO required).

### Deprecated APIs

- `security.WithAuthenticatedSecurityHeaders` → use `WithNoStoreNoReferrer`.

### Technical debt / limitations

- QR code library uses an older version (2020 commit hash).
- No built-in metrics/tracing (consumers should add their own observability).

---

## 10. How AI tools and RAG should use this file

### Consumption guidance

- **Read first** – always read `CONTEXT.md` before making architectural changes.
- **Combine with README.md** – README has usage examples; CONTEXT.md has contracts.
- **Package docs** – each `<pkg>/<pkg>.go` file has a `// Package ...` comment.

### Rules for AI agents

1. **Respect conventions** – follow naming, error handling, and testing patterns.
2. **Context-first** – all blocking functions must accept `context.Context`.
3. **No new dependencies** – prefer stdlib; justify any new vendor library.
4. **Test coverage** – new code must include `*_test.go` files.
5. **Update CONTEXT.md** – if adding a new package or pattern, update this file.

### Reference hierarchy

1. `CONTEXT.md` (architecture, conventions, contracts)
2. `README.md` (usage examples, installation)
3. Package doc comments (`<pkg>/<pkg>.go`)
4. Inline code comments (implementation details)
