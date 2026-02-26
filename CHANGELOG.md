# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-26

### Added

- `Validator` type for declaring environment variable fields
- `Field` struct with Key, Kind, Required, Default, Description, and AllowedValues
- Six supported kinds: `KindString`, `KindInteger`, `KindFloat`, `KindBoolean`, `KindURL`, `KindDuration`
- `Validate` method reads from the real process environment via `os.Getenv`
- `ValidateMap` method for deterministic testing without touching the process environment
- `Result` type with typed accessors: `String`, `Integer`, `Float`, `Boolean`, `Duration`, `Raw`
- `DurationResult` convenience function for returning `time.Duration` values
- `Schema` method returning `[]FieldSchema` — a machine-readable, JSON-safe description of all declared fields
- `ValidationError` and `ValidationErrors` structured error types
- Full `context.Context` support with cancellation
- Zero external dependencies
- Table-driven tests with race detector enabled
- GitHub Actions CI across Go 1.21, 1.22, and 1.23
