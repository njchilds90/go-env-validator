# go-env-validator

Declarative environment variable validation, type coercion, default value assignment, and machine-readable schema documentation for Go applications.

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-env-validator.svg)](https://pkg.go.dev/github.com/njchilds90/go-env-validator)
[![CI](https://github.com/njchilds90/go-env-validator/actions/workflows/ci.yml/badge.svg)](https://github.com/njchilds90/go-env-validator/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/go-env-validator)](https://goreportcard.com/report/github.com/njchilds90/go-env-validator)

## Why

Every Go application reads environment variables. Every team writes the same fragile `os.Getenv` + manual type parsing + startup panic. This library solves that with a single, zero-dependency, idiomatic declaration.

- Declare all your environment variables in one place
- Get structured, multi-error validation at startup
- Use typed accessors with no casting boilerplate
- Export a machine-readable schema for tooling and AI agents

## Installation
```bash
go get github.com/njchilds90/go-env-validator
```

## Quick Start
```go
package main

import (
    "context"
    "log"
    "time"

    envvalidator "github.com/njchilds90/go-env-validator"
)

func main() {
    v := envvalidator.New(
        envvalidator.Field{
            Key:         "PORT",
            Kind:        envvalidator.KindInteger,
            Default:     "8080",
            Description: "HTTP server listen port",
        },
        envvalidator.Field{
            Key:         "DATABASE_URL",
            Kind:        envvalidator.KindURL,
            Required:    true,
            Description: "Postgres connection URL",
        },
        envvalidator.Field{
            Key:         "LOG_LEVEL",
            Kind:        envvalidator.KindString,
            Default:     "info",
            AllowedValues: []string{"debug", "info", "warn", "error"},
            Description: "Logging verbosity",
        },
        envvalidator.Field{
            Key:         "REQUEST_TIMEOUT",
            Kind:        envvalidator.KindDuration,
            Default:     "30s",
            Description: "Timeout for outbound HTTP requests",
        },
        envvalidator.Field{
            Key:         "FEATURE_NEW_DASHBOARD",
            Kind:        envvalidator.KindBoolean,
            Default:     "false",
            Description: "Enable the new dashboard feature",
        },
    )

    result, err := v.Validate(context.Background())
    if err != nil {
        log.Fatal(err) // prints all errors at once, structured
    }

    port       := result.Integer("PORT")
    dbURL      := result.String("DATABASE_URL")
    logLevel   := result.String("LOG_LEVEL")
    timeout    := envvalidator.DurationResult(result, "REQUEST_TIMEOUT")
    newDash    := result.Boolean("FEATURE_NEW_DASHBOARD")

    _ = port
    _ = dbURL
    _ = logLevel
    _ = timeout
    _ = newDash
}
```

## Machine-Readable Schema

For tooling, documentation generators, and AI agents:
```go
schema := v.Schema()
// schema is []envvalidator.FieldSchema — safe to encode as JSON

import "encoding/json"
data, _ := json.MarshalIndent(schema, "", "  ")
fmt.Println(string(data))
```

Output:
```json
[
  {
    "key": "PORT",
    "kind": "integer",
    "required": false,
    "default": "8080",
    "description": "HTTP server listen port",
    "allowed_values": []
  },
  {
    "key": "DATABASE_URL",
    "kind": "url",
    "required": true,
    "description": "Postgres connection URL",
    "allowed_values": []
  }
]
```

## Testing Without os.Getenv

Use `ValidateMap` to test your configuration logic without touching the real environment:
```go
result, err := v.ValidateMap(context.Background(), map[string]string{
    "DATABASE_URL": "postgres://localhost/testdb",
    "PORT":         "9090",
})
```

## Supported Types

| Kind              | Accepted Input                                  | Go Type        |
|-------------------|-------------------------------------------------|----------------|
| `KindString`      | any string                                      | `string`       |
| `KindInteger`     | base-10 integer                                 | `int64`        |
| `KindFloat`       | decimal number                                  | `float64`      |
| `KindBoolean`     | true / false / 1 / 0 / yes / no (any case)     | `bool`         |
| `KindURL`         | absolute URL with scheme and host               | `string`       |
| `KindDuration`    | Go duration string: 5s, 1m30s, 2h              | `time.Duration`|

## Error Handling

If validation fails, the error is of type `ValidationErrors` (a slice of `*ValidationError`). You can inspect individual failures:
```go
result, err := v.Validate(context.Background())
if errs, ok := err.(envvalidator.ValidationErrors); ok {
    for _, e := range errs {
        fmt.Printf("field %q failed: %s\n", e.Key, e.Reason)
    }
}
```

## Philosophy

- Zero external dependencies
- No global state, no init functions, no singletons
- Pure functions where possible
- Context-aware for cancellation
- Deterministic, machine-readable output

## License

MIT
