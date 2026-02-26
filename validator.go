package envvalidator

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Validator holds a declared set of environment variable fields and provides
// methods to validate and parse values from any string-keyed map or from the
// real process environment.
type Validator struct {
	fields []Field
}

// New creates a new Validator from the given field declarations.
// Duplicate keys are not checked at construction time; the first declaration
// for a given key wins during validation.
//
// Example:
//
//	v := envvalidator.New(
//	    envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080", Description: "HTTP listen port"},
//	    envvalidator.Field{Key: "DATABASE_URL", Kind: envvalidator.KindURL, Required: true, Description: "Postgres connection URL"},
//	)
func New(fields ...Field) *Validator {
	return &Validator{fields: fields}
}

// Validate reads environment variables from the real process environment using
// os.Getenv, validates them against the declared fields, and returns a Result.
//
// If any field fails validation, a ValidationErrors value is returned. The
// caller should check for this type to inspect individual failures.
//
// Example:
//
//	result, err := v.Validate(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	port := result.Integer("PORT")
func (v *Validator) Validate(ctx context.Context) (*Result, error) {
	env := make(map[string]string)
	for _, f := range v.fields {
		if val := os.Getenv(f.Key); val != "" {
			env[f.Key] = val
		}
	}
	return v.ValidateMap(ctx, env)
}

// ValidateMap validates the given key-value map against the declared fields
// and returns a Result. This method is preferred for testing because it does
// not read from os.Getenv.
//
// Example:
//
//	result, err := v.ValidateMap(context.Background(), map[string]string{
//	    "PORT": "9090",
//	    "DATABASE_URL": "postgres://localhost/mydb",
//	})
func (v *Validator) ValidateMap(ctx context.Context, env map[string]string) (*Result, error) {
	var errs ValidationErrors
	values := make(map[string]any, len(v.fields))

	for _, f := range v.fields {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		kind := f.Kind
		if kind == "" {
			kind = KindString
		}

		raw, present := env[f.Key]
		if !present || raw == "" {
			if f.Required && f.Default == "" {
				errs = append(errs, &ValidationError{
					Key:    f.Key,
					Reason: "required variable is missing or empty",
				})
				continue
			}
			raw = f.Default
		}

		if len(f.AllowedValues) > 0 {
			found := false
			for _, allowed := range f.AllowedValues {
				if raw == allowed {
					found = true
					break
				}
			}
			if !found {
				errs = append(errs, &ValidationError{
					Key:    f.Key,
					Reason: fmt.Sprintf("value %q is not one of the allowed values: %s", raw, strings.Join(f.AllowedValues, ", ")),
				})
				continue
			}
		}

		parsed, err := parseValue(f.Key, kind, raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		values[f.Key] = parsed
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return &Result{values: values}, nil
}

// parseValue converts a raw string into the Go type corresponding to kind.
func parseValue(key string, kind Kind, raw string) (any, *ValidationError) {
	switch kind {
	case KindString:
		return raw, nil

	case KindInteger:
		n, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("cannot parse %q as an integer", raw)}
		}
		return n, nil

	case KindFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("cannot parse %q as a float", raw)}
		}
		return f, nil

	case KindBoolean:
		normalized := strings.ToLower(strings.TrimSpace(raw))
		switch normalized {
		case "true", "1", "yes":
			return true, nil
		case "false", "0", "no":
			return false, nil
		default:
			return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("cannot parse %q as a boolean; accepted values are true, false, 1, 0, yes, no", raw)}
		}

	case KindURL:
		trimmed := strings.TrimSpace(raw)
		u, err := url.ParseRequestURI(trimmed)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("cannot parse %q as an absolute URL with scheme and host", raw)}
		}
		return trimmed, nil

	case KindDuration:
		d, err := time.ParseDuration(strings.TrimSpace(raw))
		if err != nil {
			return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("cannot parse %q as a duration; use Go duration syntax such as 5s, 1m30s, or 2h", raw)}
		}
		return d, nil

	default:
		return nil, &ValidationError{Key: key, Reason: fmt.Sprintf("unknown kind %q", kind)}
	}
}

// DurationResult is a convenience wrapper that returns a time.Duration from a
// Result. It panics if the key was not declared or is not a KindDuration field.
//
// Example:
//
//	timeout := envvalidator.DurationResult(result, "TIMEOUT")
func DurationResult(r *Result, key string) time.Duration {
	v, ok := r.Raw(key)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q was not declared in the validator", key))
	}
	d, ok := v.(time.Duration)
	if !ok {
		panic(fmt.Sprintf("env-validator: key %q is not a duration field", key))
	}
	return d
}
