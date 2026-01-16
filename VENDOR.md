# VENDOR.md

## Overview

Cloud Native Utils maintains a small, focused set of external dependencies. Each vendor library addresses a specific cross-cutting concern that would be impractical to implement from scratch. The project philosophy favors:

- **Minimal dependencies** – only what is genuinely needed
- **Stable, well-maintained libraries** – from reputable sources
- **Pure Go when possible** – avoiding CGO for portability

All direct dependencies are declared in `go.mod`. This document explains when and how to use each one.

---

## Approved Vendor Libraries

### github.com/coreos/go-oidc/v3

- **Purpose**: OpenID Connect (OIDC) identity provider integration for authentication flows.
- **Repository**: https://github.com/coreos/go-oidc
- **Version**: v3.17.0

**Key Packages/Modules**:
- `oidc` – Provider discovery, ID token verification, claims extraction

**When to use it**:
- Implementing OIDC-based authentication (Google, Azure AD, Keycloak, etc.)
- Verifying ID tokens from identity providers
- Extracting user claims (email, name, subject) from tokens

**Integration pattern**:
- Used exclusively in `security/identity_provider.go`
- Always pair with `golang.org/x/oauth2` for the OAuth2 flow
- Use `oidc.NewProvider()` for provider discovery
- Use `provider.Verifier()` to create a token verifier

**Example**:
```go
import "github.com/coreos/go-oidc/v3/oidc"

provider, _ := oidc.NewProvider(ctx, issuerURL)
verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
idToken, _ := verifier.Verify(ctx, rawToken)

var claims IdentityTokenClaims
_ = idToken.Claims(&claims)
```

**Cautions**:
- Requires a valid OIDC provider URL at runtime
- Network-dependent: provider discovery makes HTTP calls
- Configure via environment variables: `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_REDIRECT_URL`
- Session IDs are stored in HTTP-only secure cookies (`sid`), not URL paths
- Logout reads session ID from cookie and clears it with `MaxAge: -1`

---

### golang.org/x/oauth2

- **Purpose**: OAuth2 client flows for authorization code exchange and token management.
- **Repository**: https://github.com/golang/oauth2
- **Version**: v0.34.0

**Key Packages/Modules**:
- Root package – `oauth2.Config`, token exchange, PKCE support

**When to use it**:
- Building OAuth2 authorization code flows
- Exchanging authorization codes for tokens
- PKCE (Proof Key for Code Exchange) flows

**Integration pattern**:
- Used in `security/identity_provider.go` alongside `go-oidc`
- Configure `oauth2.Config` with endpoints from OIDC provider
- Use `oauth2.SetAuthURLParam()` for PKCE code verifier

**Example**:
```go
import "golang.org/x/oauth2"

config := &oauth2.Config{
    ClientID:     clientID,
    ClientSecret: clientSecret,
    RedirectURL:  redirectURL,
    Endpoint:     provider.Endpoint(),
    Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
}

// Exchange code with PKCE
token, _ := config.Exchange(ctx, code,
    oauth2.SetAuthURLParam("code_verifier", codeVerifier),
)
```

**Cautions**:
- Store client secrets securely (environment variables, not code)
- Use PKCE for public clients (SPAs, mobile apps)
- Token expiration must be handled by the application

---

### golang.org/x/crypto

- **Purpose**: Cryptographic primitives, specifically bcrypt for password hashing.
- **Repository**: https://github.com/golang/crypto
- **Version**: v0.46.0

**Key Packages/Modules**:
- `bcrypt` – Password hashing and verification

**When to use it**:
- Hashing user passwords before storage
- Verifying passwords during authentication
- **Do not use** for encrypting data (use AES-GCM from `crypto/aes` instead)

**Integration pattern**:
- Used in `security/password.go`
- Wrapped by `security.Password()` and `security.IsPasswordValid()`
- **Prefer the wrapper functions** over direct bcrypt calls

**Example**:
```go
import "github.com/andygeiss/cloud-native-utils/security"

// Hash a password (cost 14)
hash, err := security.Password([]byte("p@ssw0rd"))

// Verify a password
ok := security.IsPasswordValid(hash, []byte("p@ssw0rd"))
```

**Cautions**:
- bcrypt is intentionally slow; cost factor 14 is used (secure default)
- Never store plaintext passwords
- bcrypt output includes salt; no separate salt storage needed

---

### github.com/segmentio/kafka-go

- **Purpose**: Apache Kafka client for distributed messaging.
- **Repository**: https://github.com/segmentio/kafka-go
- **Version**: v0.4.49

**Key Packages/Modules**:
- Root package – `kafka.Writer`, `kafka.Reader`, `kafka.Message`

**When to use it**:
- Publishing messages to Kafka topics
- Subscribing to Kafka topics for event-driven processing
- Building distributed, decoupled microservice communication

**Integration pattern**:
- Used exclusively in `messaging/dispatcher_external.go`
- Access via `messaging.NewExternalDispatcher()` – do not use kafka-go directly
- Configure brokers via `KAFKA_BROKERS` environment variable (comma-separated)
- Stability patterns (retry, timeout) are applied automatically

**Example**:
```go
import "github.com/andygeiss/cloud-native-utils/messaging"

dispatcher := messaging.NewExternalDispatcher()

// Publish
_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", payload))

// Subscribe
_ = dispatcher.Subscribe(ctx, "user.created", func(ctx context.Context, msg messaging.Message) (messaging.MessageState, error) {
    // Handle message
    return messaging.MessageStateProcessed, nil
})
```

**Cautions**:
- Requires running Kafka cluster
- Set `KAFKA_BROKERS` environment variable (e.g., `localhost:9092,localhost:9093`)
- Auto topic creation is enabled; manage topics explicitly in production
- For local development without Kafka, use `messaging.NewInternalDispatcher()` instead

---

### github.com/skip2/go-qrcode

- **Purpose**: QR code generation as PNG images or data URLs.
- **Repository**: https://github.com/skip2/go-qrcode
- **Version**: v0.0.0-20200617195104-da1b6568686e

**Key Packages/Modules**:
- Root package – `qrcode.Encode()`, recovery levels

**When to use it**:
- Generating QR codes for URLs, text, or data
- Embedding QR codes in HTML via data URLs
- TOTP/2FA setup screens

**Integration pattern**:
- Wrapped by `imaging.QRCodeGenerator`
- **Prefer the wrapper** for consistent sizing and recovery levels

**Example**:
```go
import "github.com/andygeiss/cloud-native-utils/imaging"

gen := imaging.NewQRCodeGenerator().
    WithSize(256).
    WithRecoveryLevel(imaging.RecoveryMedium)

// Get as data URL for HTML img src
dataURL, _ := gen.DataURL("https://example.com")

// Get as raw PNG bytes
png, _ := gen.PNG("https://example.com")
```

**Cautions**:
- Larger data = larger QR code; keep content concise
- Higher recovery levels increase QR code density
- Library is archived but stable; no active development

---

### gopkg.in/yaml.v3

- **Purpose**: YAML parsing and serialization for configuration and data files.
- **Repository**: https://github.com/go-yaml/yaml
- **Version**: v3.0.1

**Key Packages/Modules**:
- Root package – `yaml.Marshal()`, `yaml.Unmarshal()`

**When to use it**:
- YAML-based resource persistence (`resource.YamlFileAccess`)
- Configuration files (when JSON is too verbose)

**Integration pattern**:
- Used in `resource/yaml_file_access.go` for YAML-based CRUD storage
- Always use `yaml.v3` (not v2) for improved performance and features

**Example**:
```go
import "gopkg.in/yaml.v3"

// Unmarshal
var config map[string]any
_ = yaml.Unmarshal(data, &config)

// Marshal
data, _ := yaml.Marshal(config)
```

**Cautions**:
- YAML is whitespace-sensitive; be careful with indentation
- Use struct tags (`yaml:"field_name"`) for custom field mapping
- For simple key-value data, JSON may be more portable

---

### modernc.org/sqlite

- **Purpose**: Pure-Go SQLite database driver (no CGO required).
- **Repository**: https://gitlab.com/cznic/sqlite
- **Version**: v1.40.1

**Key Packages/Modules**:
- Root package – imported as `_ "modernc.org/sqlite"` for driver registration

**When to use it**:
- Local SQLite database access without CGO
- Embedded databases for development or edge deployments
- When `database/sql` with SQLite is needed

**Integration pattern**:
- Used in `resource/sqlite_access.go` via standard `database/sql` interface
- Import as blank identifier for side-effect registration: `_ "modernc.org/sqlite"`
- Access via `resource.NewSqliteAccess[K, V](db)` wrapper

**Example**:
```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

db, _ := sql.Open("sqlite", "file:data.db?mode=rwc")
defer db.Close()

store := resource.NewSqliteAccess[string, User](db)
_ = store.Init(ctx) // Creates kv_store table
_ = store.Create(ctx, "user-1", user)
```

**Cautions**:
- Pure Go = slower than CGO-based `mattn/go-sqlite3` but more portable
- Uses a fixed table schema (`kv_store` with `key TEXT, value TEXT`)
- JSON-encodes values; not suitable for complex SQL queries
- Call `Init()` to create the required table structure

---

## Cross-cutting Concerns and Recommended Patterns

### Authentication & Authorization

| Concern | Recommended Vendor | Notes |
|---------|-------------------|-------|
| OIDC authentication | `go-oidc/v3` + `oauth2` | Use together for complete flow |
| Password hashing | `x/crypto/bcrypt` | Via `security.Password()` wrapper |

### Data Persistence

| Concern | Recommended Vendor | Notes |
|---------|-------------------|-------|
| YAML files | `yaml.v3` | Via `resource.YamlFileAccess` |
| SQLite | `modernc.org/sqlite` | Via `resource.SqliteAccess` |
| JSON files | Standard library | Via `resource.JsonFileAccess` |
| In-memory | Standard library | Via `resource.InMemoryAccess` |

### Messaging

| Concern | Recommended Vendor | Notes |
|---------|-------------------|-------|
| Kafka messaging | `kafka-go` | Via `messaging.NewExternalDispatcher()` |
| In-memory pub/sub | Standard library | Via `messaging.NewInternalDispatcher()` |

### Imaging

| Concern | Recommended Vendor | Notes |
|---------|-------------------|-------|
| QR codes | `go-qrcode` | Via `imaging.QRCodeGenerator` |

---

## Indirect Dependencies

The following are transitive dependencies pulled in by direct dependencies. Do not import these directly:

| Dependency | Pulled by | Purpose |
|------------|-----------|---------|
| `github.com/go-jose/go-jose/v4` | `go-oidc` | JWT/JWS handling |
| `github.com/klauspost/compress` | `kafka-go` | Compression |
| `github.com/pierrec/lz4/v4` | `kafka-go` | LZ4 compression |
| `github.com/google/uuid` | `sqlite` | UUID generation |
| `modernc.org/libc` | `sqlite` | C runtime emulation |

---

## Vendors to Avoid

The following are explicitly **not recommended** for this project:

| Vendor | Reason | Use Instead |
|--------|--------|-------------|
| `github.com/mattn/go-sqlite3` | Requires CGO | `modernc.org/sqlite` |
| `github.com/stretchr/testify` | External test dependency | `assert.That()` from this repo |
| `github.com/golang/mock` | External test dependency | `resource.MockAccess` or manual mocks |
| `github.com/sirupsen/logrus` | Superseded by slog | `log/slog` via `logging.NewJsonLogger()` |
| `gopkg.in/yaml.v2` | Outdated | `gopkg.in/yaml.v3` |

---

## Version Policy

- **Pin exact versions** in `go.mod` (Go modules default behavior)
- **Update dependencies** via `go get -u <package>@latest` and run tests
- **Security updates** should be applied promptly; check `go list -m -u all`
- **Breaking changes** should be documented in commit messages

---

## Adding New Dependencies

Before adding a new dependency:

1. **Check if existing vendors cover the use case** – prefer wrappers over new libraries
2. **Evaluate alternatives** – prefer standard library or well-maintained options
3. **Avoid CGO** – pure Go is preferred for portability
4. **Document in VENDOR.md** – add a section following the template above
5. **Create a wrapper** – expose vendor functionality through project-specific APIs
