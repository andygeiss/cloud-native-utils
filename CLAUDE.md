# Cloud Native Utils

## Documentation Policy

### Update CLAUDE.md when changes affect:
1. **Architecture** - New patterns, packages, dependencies
2. **API surface** - New interfaces, functions, types
3. **Conventions** - New naming rules, anti-patterns, gotchas
4. **Decisions** - Architectural trade-offs or technology choices

### Update README.md when changes affect:
1. **User-facing behavior** - New packages, usage examples
2. **Setup instructions** - New prerequisites, Go version requirements

### Documentation Checklist (before commit):
- [ ] New package → Add to Project Structure and README?
- [ ] New pattern/gotcha → Add to Conventions or Gotchas section?
- [ ] New interface → Document in Package Patterns?
- [ ] New environment variable → Add to `.env.example`?
- [ ] Roadmap item completed → Mark as [x]?

### Documentation Update Matrix

| Change Type | CLAUDE.md Section | README.md |
|-------------|-------------------|-----------|
| New package | Project Structure, Roadmap | Features, Usage |
| New interface | Package Patterns | Usage examples |
| New convention | Coding Conventions | - |
| New gotcha | Gotchas | - |
| Architectural decision | Decisions | - |
| New environment variable | - | .env.example |

---

## Project Structure

```
cloud-native-utils/
├── assert/          Test assertions (assert.That)
├── consistency/     Transactional event log (JSON file persistence)
├── efficiency/      Channel helpers, gzip middleware, similarity search, sparse data structures
├── env/             Generic environment variable parsing
├── event/           Domain event interfaces
├── extensibility/   Dynamic Go plugin loading
├── logging/         Structured JSON logging (log/slog)
├── mcp/             Model Context Protocol server (Claude Desktop)
├── messaging/       Pub-sub dispatchers (in-memory, Kafka)
├── resource/        Generic CRUD backends (memory, sharded-sparse, JSON, YAML, SQLite, PostgreSQL)
├── security/        AES-GCM encryption, password hashing, HMAC
├── service/         Context helpers, lifecycle management
├── slices/          Generic slice utilities
├── stability/       Resilience patterns (breaker, retry, throttle, debounce, timeout)
├── templating/      HTML template engine (embed.FS support)
└── web/             HTTP server, client, sessions, OIDC, auth middleware
```

### File Naming Conventions

| Pattern | Purpose |
|---------|---------|
| `{feature}.go` | Main implementation |
| `{feature}_test.go` | Tests for feature |
| `{package}.go` | Package-level documentation |
| `access.go` | Interface definitions |
| `testdata/` | Test fixtures and data |

---

## Commands

```
just test            Run all tests with coverage
just test-integration Run integration tests (requires tags)
just benchmark       Run consistency benchmarks
just lint            Run golangci-lint
just plugin          Build test plugins
just make-certs      Generate mTLS certificates
```

---

## Coding Conventions

### Generic Type Patterns

```go
// CRUD access with comparable key and any value
type Access[K comparable, V any] interface { ... }

// Implementation with generics
type InMemoryAccess[K comparable, V any] struct { ... }
```

### Function Signature Convention

```go
// Cloud-native service function pattern
type Function[IN, OUT any] func(ctx context.Context, in IN) (out OUT, err error)
```

### Constructor Naming

```go
// Always use New{Type} pattern
func NewInMemoryAccess[K comparable, V any]() *InMemoryAccess[K, V]
func NewServer(handler http.Handler) *Server
func NewJsonLogger() *slog.Logger
```

### Error Handling

```go
// Use named error constants
var (
    ErrorResourceAlreadyExists = errors.New("resource already exists")
    ErrorResourceNotFound      = errors.New("resource not found")
)

// Prefix with Err for sentinel errors
var ErrBreakerServiceUnavailable = errors.New("service unavailable")
```

### Context Propagation

- All operations require `context.Context` as first parameter
- Check context cancellation early in operations
- Propagate context through all function calls

### Concurrency Safety

```go
// Use sync.RWMutex for shared state
type InMemoryAccess[K comparable, V any] struct {
    data map[K]V
    mu   sync.RWMutex
}

// Read lock for read operations
func (a *InMemoryAccess[K, V]) Read(ctx context.Context, key K) (*V, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    // ...
}

// Write lock for mutations
func (a *InMemoryAccess[K, V]) Create(ctx context.Context, key K, value V) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    // ...
}
```

---

## Testing Conventions

### AAA Pattern

```go
func Test_Feature_With_Condition_Should_Outcome(t *testing.T) {
    // Arrange
    access := resource.NewInMemoryAccess[string, User]()
    ctx := context.Background()
    user := User{Name: "Alice"}

    // Act
    err := access.Create(ctx, "user-1", user)

    // Assert
    assert.That(t, "error should be nil", err, nil)
}
```

### Test Naming

```
Test_{Feature}_With_{Condition}_Should_{Outcome}

Examples:
- Test_Breaker_With_ThresholdExceeded_Should_ReturnServiceUnavailable
- Test_InMemoryAccess_With_DuplicateKey_Should_ReturnError
- Test_Retry_With_ContextCanceled_Should_ReturnContextError
```

### Assertion Usage

```go
// Use assert.That for all assertions
assert.That(t, "description of what is being tested", got, expected)

// Examples
assert.That(t, "result should be 42", result, 42)
assert.That(t, "error should be nil", err, nil)
assert.That(t, "user name should match", user.Name, "Alice")
```

### Test Scenarios to Cover

1. Happy path
2. Context cancellation
3. Context timeout
4. Error conditions
5. Concurrent access
6. Recovery/resilience (for stability patterns)

---

## Environment Variables

See `.env.example` for the full list. Key variables:

| Variable | Package | Description |
|----------|---------|-------------|
| `PORT` | web | HTTP server port |
| `SERVER_READ_TIMEOUT` | web | HTTP read timeout |
| `SERVER_WRITE_TIMEOUT` | web | HTTP write timeout |
| `KAFKA_BROKERS` | messaging | Kafka broker addresses |
| `OIDC_ISSUER_URL` | web | OpenID Connect issuer |
| `OIDC_CLIENT_ID` | web | OIDC client ID |
| `OIDC_CLIENT_SECRET` | web | OIDC client secret |

---

## Roadmap

### Completed
- [x] Core packages (assert, env, logging, service, slices)
- [x] Stability patterns (breaker, retry, throttle, debounce, timeout)
- [x] Resource backends (in-memory, JSON, YAML, SQLite, PostgreSQL)
- [x] Messaging (internal dispatcher, Kafka dispatcher)
- [x] Event interfaces (Event, EventPublisher, EventSubscriber)
- [x] Security (AES-GCM, password hashing, HMAC)
- [x] Web (server, client, sessions, OIDC)
- [x] MCP server (Model Context Protocol for AI tools)
- [x] Bearer token authentication middleware
- [x] ShardedSparseAccess for high-concurrency workloads
- [x] Similarity search (Cosine, Jaccard) for sparse vector/set data

### Planned
- [ ] Redis backend for resource package
- [ ] Metrics/observability package
- [ ] Rate limiting middleware

---

## Decisions

| Decision | Rationale |
|----------|-----------|
| Single-responsibility packages | Each package solves one problem, import only what you need |
| Generic CRUD interface | Swap backends without changing domain code |
| `Function[IN, OUT]` signature | Composable with stability wrappers |
| Context-first parameters | Cloud-native pattern, cancellation support |
| No global state | Testability, concurrent safety |
| `sync.RWMutex` over channels | Simpler for CRUD operations |
| golangci-lint | Consistent code quality across packages |
| Sharding + sparse-dense for high-perf storage | 3-4x concurrent throughput, O(1) delete, cache-friendly iteration |
| KeyedSparseSet → SparseSharding → ShardedSparseAccess | Layered composition: data structure, concurrency, CRUD semantics |

---

## Package Patterns

### Generic CRUD (resource)

```go
// Single interface, multiple implementations
type Access[K, V any] interface {
    Create(ctx context.Context, key K, value V) error
    Read(ctx context.Context, key K) (*V, error)
    Update(ctx context.Context, key K, value V) error
    Delete(ctx context.Context, key K) error
    List(ctx context.Context) ([]V, error)
}

// Implementations
store := resource.NewInMemoryAccess[string, User]()
store := resource.NewShardedSparseAccess[string, User](32)  // High-performance: 32 shards
store := resource.NewJsonFileAccess[string, User]("users.json")
store := resource.NewPostgresAccess[string, User](db)
```

### Resilience Wrappers (stability)

```go
// Wrap any Function[IN, OUT] with resilience patterns
fn := stability.Breaker(yourFunc, 3)          // Opens after 3 failures
fn := stability.Retry(fn, 5, time.Second)     // Retry 5 times
fn := stability.Throttle(fn, 10)              // Max 10 concurrent
fn := stability.Timeout(fn, 5*time.Second)    // 5s timeout
```

### Event-Driven (event + messaging)

```go
// Domain event interface
type Event interface {
    Topic() string
}

// Publish/subscribe with dispatcher
dispatcher := messaging.NewInternalDispatcher()
_ = dispatcher.Subscribe(ctx, "user.created", handler)
_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", payload))
```

### HTTP with Auth (web)

```go
// Session-based auth for web UI
mux.HandleFunc("GET /protected", web.WithAuth(sessions, handler))

// Bearer token auth for APIs/MCP
mux.HandleFunc("POST /mcp", web.WithBearerAuth(verifier, handler))
```

### Sparse Data Structures (efficiency)

```go
// KeyedSparseSet: O(1) operations, cache-friendly iteration
// Uses bidirectional key mapping (sparse-dense) for O(1) delete via swap-remove
set := efficiency.NewKeyedSparseSet[string, User](100)
isNew := set.Put("user-1", user)  // Returns true if key was new
value := set.Get("user-1")         // Returns *User or nil
deleted := set.Delete("user-1")    // O(1) swap-remove

// SparseSharding: Concurrent access with per-shard locking
// Wraps KeyedSparseSet with FNV-1a hash distribution
shards := efficiency.NewSparseSharding[string, User](32)
shards.Put("user-1", user)
shards.ForEach(func(k string, v User) bool { return true })
shards.ForEachShard(func(idx int, iterate func(fn func(string, User) bool)) {
    // Per-shard iteration with cancellation checks between shards
    iterate(func(k string, v User) bool { return true })
})
```

### Similarity Search (resource + efficiency)

```go
// SearchSimilar is a method on ShardedSparseAccess
// Use a custom scorer function with utility functions from efficiency package

// Cosine similarity for TF-IDF vectors
results := store.SearchSimilar(ctx, func(doc Document) float64 {
    return efficiency.CosineSimilarity(
        query.Indices, doc.Indices,
        query.Values, doc.Values,
        query.Norm, doc.Norm,
    )
}, resource.SearchOptions{TopK: 10, Threshold: 0.5})

// Jaccard similarity for tag/term sets
results := store.SearchSimilar(ctx, func(article Article) float64 {
    return efficiency.JaccardSimilarity(queryTags, article.Tags)
}, resource.SearchOptions{TopK: 10})

// Custom scoring (e.g., weighted combination)
results := store.SearchSimilar(ctx, func(item Item) float64 {
    return 0.7*textScore(query, item) + 0.3*tagScore(query, item)
}, resource.SearchOptions{TopK: 5})
```

---

## Gotchas

1. **Context is mandatory** - All operations require `context.Context` as the first parameter. Never pass `nil`.

2. **Generic constraints** - Use `comparable` constraint for map keys:
   ```go
   // Correct
   type Access[K comparable, V any] interface { ... }

   // Wrong - K must be comparable for map usage
   type Access[K any, V any] interface { ... }
   ```

3. **PostgreSQL Init required** - Call `Init(ctx)` before using PostgresAccess to create tables:
   ```go
   store := resource.NewPostgresAccess[string, User](db)
   _ = store.Init(ctx) // Creates kv_store table
   ```

4. **Circuit breaker state** - Breaker uses exponential backoff for recovery. Test with sufficient wait time.

5. **Kafka dispatcher** - Requires `KAFKA_BROKERS` environment variable. Use `NewInternalDispatcher()` for tests.

6. **Test package separation** - Tests use `package {name}_test` to test public API only.

7. **Mutex ordering** - Always acquire locks in consistent order to avoid deadlocks. Use RLock for reads, Lock for writes.

8. **ShardedSparseAccess memory trade-off** - Uses ~2x memory vs InMemoryAccess due to bidirectional key mapping. Use when concurrent throughput is critical; use InMemoryAccess for memory-constrained scenarios.

9. **Similarity search requires sorted indices** - `CosineSimilarity` and `JaccardSimilarity` utility functions require index slices to be sorted in ascending order for O(m+n) merge-loop efficiency. Pre-compute and cache norms for cosine similarity.
