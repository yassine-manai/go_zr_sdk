# Package Architecture Recommendations

## Overall Structure

**Root Package Name**: Choose something descriptive like `thirdparty-client` or `<vendor>-sdk`

```
project-root/
├── client/          # Main client initialization
├── ui/              # UI service interactions
├── db/              # Database operations
├── config/          # Configuration management
├── models/          # Shared data structures
├── errors/          # Custom error types
├── logger/          # Logging abstraction
├── internal/        # Private utilities
└── examples/        # Usage examples
```

---

## Detailed Breakdown

### 1. Folder Organization

**UI Services** (organize by domain):
```
ui/
├── customer_media/
│   ├── service.go       # Service interface & implementation
│   ├── models.go        # Domain-specific models
│   └── validators.go    # Input validation
├── ticket_class/
├── rebate/
├── shift/
└── common.go           # Shared UI utilities
```

**DB Layer**:
```
db/
├── queries/
│   ├── customer_media.sql    # SQL files or constants
│   ├── ticket_class.sql
│   └── ...
├── repository.go             # DB interface
├── postgres.go / mysql.go    # DB-specific implementations
└── transaction.go            # Transaction helpers
```

---

### 2. Architecture Patterns

**Client Pattern** (recommended):
- Single entry point: `client.New(config)` returns a client with all services
- Each service (UI/DB) is accessible via the client
- Client manages shared concerns (auth, retry, connection pooling)

**Layering**:
1. **Transport Layer**: HTTP client, DB connection handling
2. **Service Layer**: Business logic per domain (customer_media, etc.)
3. **Data Layer**: Models, serialization/deserialization
4. **Error Layer**: Standardized error handling

**Key Design Principles**:
- **Interface-based**: Define interfaces for each service for testability
- **Context-aware**: All methods accept `context.Context` for cancellation/timeouts
- **Immutable config**: Configuration set once at initialization
- **Thread-safe**: Client should be safe for concurrent use

---

### 3. Logging Strategy

**Structured Logging Approach**:

**Logger Interface** (abstraction):
- Allow users to plug in their own logger (logrus, zap, slog)
- Provide default no-op or simple logger
- Use interface: `Logger` with methods: Debug, Info, Warn, Error

**What to Log**:
- **Request level**: URL, method, headers (sanitized), request ID
- **Response level**: Status code, duration, response size
- **Error level**: Full error context, stack traces
- **DB queries**: Query text (parameterized), execution time, rows affected

**Log Levels**:
- DEBUG: Full request/response bodies (disabled in production)
- INFO: Request summary, successful operations
- WARN: Retries, deprecated feature usage
- ERROR: Failures, validation errors

**Best Practices**:
- Add correlation IDs for tracing requests across services
- Include timestamps, service name, version
- Never log sensitive data (passwords, tokens) - sanitize
- Make logging configurable per service/operation

---

### 4. Configuration Management

**Config Structure**:
- **Connection settings**: Base URLs, timeouts, retry policies
- **Authentication**: API keys, tokens, credentials
- **Feature flags**: Enable/disable specific services
- **Logging**: Log level, output destination

**Config Sources** (priority order):
1. Code (programmatic config)
2. Environment variables
3. Config files (YAML/JSON)
4. Defaults

**Validation**:
- Validate all config at initialization time
- Fail fast with clear error messages
- Provide sensible defaults

---

### 5. Error Handling

**Custom Error Types**:
- `ValidationError`: Input validation failures
- `AuthenticationError`: Auth failures
- `RateLimitError`: Rate limiting
- `ServiceUnavailableError`: 3rd party down
- `NetworkError`: Connection issues

**Error Wrapping**:
- Use Go 1.13+ error wrapping (`fmt.Errorf` with `%w`)
- Preserve original error context
- Add package-specific context

**Retry Logic**:
- Implement exponential backoff with jitter
- Make retryable errors configurable
- Expose retry configuration to users

---

### 6. Additional Components

**Authentication/Authorization**:
- Centralized auth handling in the client
- Support multiple auth methods (API key, OAuth, JWT)
- Auto-refresh tokens if applicable

**Rate Limiting**:
- Client-side rate limiting to respect 3rd party limits
- Queue or reject requests when limit reached

**Connection Pooling**:
- For DB: Configure connection pool settings
- For HTTP: Reuse HTTP clients, connection pooling

**Middleware/Interceptors**:
- Request/response logging
- Metrics collection
- Request ID injection
- Authentication injection

**Testing Support**:
- Provide mock implementations
- Test helpers for common scenarios
- Integration test utilities

---

### 7. Scope & Features

**Core Features** (MVP):
- ✅ Client initialization with config
- ✅ All UI service methods
- ✅ All DB query methods
- ✅ Error handling
- ✅ Basic logging
- ✅ Context support

**Nice-to-Have**:
- Circuit breaker pattern
- Metrics/observability (Prometheus)
- Request caching
- Pagination helpers
- Webhooks support (if 3rd party has them)
- Async/batch operations
- Graceful degradation

**Documentation**:
- godoc comments on all public APIs
- README with quick start
- Examples for each service
- Migration guide if replacing existing code
- Architecture decision records (ADRs)

---

### 8. Versioning & Release Strategy

- Use semantic versioning (v1.0.0)
- Go modules for dependency management
- Maintain CHANGELOG
- Tag releases in Git
- Consider backwards compatibility carefully
- Deprecation policy for breaking changes

---

### 9. Performance Considerations

- Connection reuse (HTTP keep-alive, DB pooling)
- Request batching where possible
- Lazy initialization of services
- Avoid unnecessary allocations
- Benchmark critical paths
- Profile memory usage

---

This structure will give developers a clean, maintainable SDK that's easy to use and extend. The key is separation of concerns, clear interfaces, and robust error handling.