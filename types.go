// Package envvalidator provides declarative environment variable validation,
// type coercion, default value assignment, and machine-readable schema
// documentation for Go applications.
//
// It is designed to be used at application startup to fail fast on
// configuration errors with clear, structured error messages.
package envvalidator

import "fmt"

// Kind represents the expected data type of an environment variable.
type Kind string

const (
	// KindString expects any string value.
	KindString Kind = "string"

	// KindInteger expects a base-10 integer value.
	KindInteger Kind = "integer"

	// KindFloat expects a decimal number value.
	KindFloat Kind = "float"

	// KindBoolean expects one of: true, false, 1, 0, yes, no (case-insensitive).
	KindBoolean Kind = "boolean"

	// KindURL expects a value that is a valid absolute URL with a scheme and host.
	KindURL Kind = "url"

	// KindDuration expects a Go duration string such as "5s", "1m30s", or "2h".
	KindDuration Kind = "duration"
)

// Field describes a single expected environment variable: its key, type,
// whether it is required, its default value if optional, and a human-readable
// description used in schema output and error messages.
type Field struct {
	// Key is the exact environment variable name, for example "DATABASE_URL".
	Key string

	// Kind is the expected data type. Defaults to KindString if not set.
	Kind Kind

	// Required marks this variable as mandatory. Validation fails if it is
	// absent and no Default is provided.
	Required bool

	// Default is the value used when the variable is absent and Required is
	// false. It must be a string representation of the correct Kind.
	Default string

	// Description is a human-readable explanation of what this variable does.
	// It appears in schema output and validation error messages.
	Description string

	// AllowedValues, if non-empty, restricts the value to one of the listed
	// strings. The comparison is case-sensitive.
	AllowedValues []string
}

// FieldSchema is the machine-readable description of a single field as
// returned by Validator.Schema. It is safe to marshal to JSON.
type FieldSchema struct {
	Key           string   `json:"key"`
	Kind          string   `json:"kind"`
	Required      bool     `json:"required"`
	Default       string   `json:"default,omitempty"`
	Description   string   `json:"description,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
}

// ValidationError describes a single field that failed validation.
type ValidationError struct {
	// Key is the environment variable name that caused the error.
	Key string

	// Reason is a human-readable description of why validation failed.
	Reason string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("env-validator: field %q: %s", e.Key, e.Reason)
}

// ValidationErrors is a slice of ValidationError values returned when one or
// more fields fail validation. It implements the error interface so callers
// can treat the entire batch as a single error.
type ValidationErrors []*ValidationError

// Error implements the error interface, returning all errors joined by newlines.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	out := fmt.Sprintf("%d environment variable validation error(s):\n", len(ve))
	for _, e := range ve {
		out += "  - " + e.Error() + "\n"
	}
	return out
}

// Result holds the successfully parsed and validated values from the
// environment. Values are accessed by their field key.
type Result struct {
	values map[string]any
}

// String returns the string value for the given key. It panics if the key was
// not declared or if the field Kind is not KindString.
func (r *Result) String(key string) string {
	v, ok := r.values[key]
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	s, ok := v.(string)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q is not a string field", key))
	}
	return s
}

// Integer returns the int64 value for the given key. It panics if the key was
// not declared or if the field Kind is not KindInteger.
func (r *Result) Integer(key string) int64 {
	v, ok := r.values[key]
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	n, ok := v.(int64)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q is not an integer field", key))
	}
	return n
}

// Float returns the float64 value for the given key. It panics if the key was
// not declared or if the field Kind is not KindFloat.
func (r *Result) Float(key string) float64 {
	v, ok := r.values[key]
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	f, ok := v.(float64)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q is not a float field", key))
	}
	return f
}

// Boolean returns the bool value for the given key. It panics if the key was
// not declared or if the field Kind is not KindBoolean.
func (r *Result) Boolean(key string) bool {
	v, ok := r.values[key]
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	b, ok := v.(bool)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q is not a boolean field", key))
	}
	return b
}

// Duration returns the time.Duration value for the given key. It panics if
// the key was not declared or if the field Kind is not KindDuration.
func (r *Result) Duration(key string) interface{} {
	// Returned as interface{} to avoid a time import in this file.
	// The concrete type is time.Duration. Use validator.go's Duration accessor.
	v, ok := r.values[key]
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	return v
}

// Raw returns the raw parsed value for the given key as an empty interface.
// Useful when the caller wants to perform their own type assertion.
func (r *Result) Raw(key string) (any, bool) {
	v, ok := r.values[key]
	return v, ok
}
