# go-common

Shared library for Go projects.

## Packages

### `env`
Loads and parses environment variables into a struct using struct tags.

### `exception`
Standardized error types and handling for consistent error responses across services.

### `format`
Formatters for common data types:
- Phone number normalization

### `hash`
Hashing utilities:
- MD5, SHA256
- Bcrypt hash and compare

### `http`
HTTP utilities:
- Server setup and configuration
- Middleware (logging, recovery, etc)
- Standardized HTTP responder

### `json`
Thin wrapper over `encoding/json`:
- Marshal, Unmarshal
- MustMarshal for panic-on-error cases

### `jwt`
JWT utilities:
- Token generation with custom claims
- Token parsing and validation

### `log`
Structured logger with dual output:
- Colored console output
- OpenTelemetry log export

### `postgres`
PostgreSQL utilities:
- Connection setup
- Query builder helpers

### `smtp`
Email sending via SMTP:
- Simple send with subject and body
- HTML email support

### `telemetry`
OpenTelemetry setup:
- Tracer and logger provider initialization
- OTLP exporter configuration

### `test`
Testing utilities:
- Common test helpers and assertions

### `validator`
Request validation wrapper over `go-playground/validator`:
- Struct validation
- Custom error messages

## Installation

```bash
go get github.com/dani-susanto/go-common@latest
```

## Usage

```go
import (
    "github.com/dani-susanto/go-common/env"
    "github.com/dani-susanto/go-common/log"
    "github.com/dani-susanto/go-common/validator"
    "github.com/dani-susanto/go-common/http"
    json "github.com/dani-susanto/go-common/json"
)
```