package envvalidator_test

import (
	"context"
	"strings"
	"testing"
	"time"

	envvalidator "github.com/njchilds90/go-env-validator"
)

func TestValidateMap_RequiredFieldMissing(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "DATABASE_URL", Kind: envvalidator.KindURL, Required: true},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required field, got nil")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("expected error to mention DATABASE_URL, got: %s", err.Error())
	}
}

func TestValidateMap_DefaultApplied(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080"},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Integer("PORT") != 8080 {
		t.Errorf("expected 8080, got %d", result.Integer("PORT"))
	}
}

func TestValidateMap_ProvidedValueOverridesDefault(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080"},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{"PORT": "9090"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Integer("PORT") != 9090 {
		t.Errorf("expected 9090, got %d", result.Integer("PORT"))
	}
}

func TestValidateMap_StringKind(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "APP_NAME", Kind: envvalidator.KindString, Default: "myapp"},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{"APP_NAME": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String("APP_NAME") != "hello" {
		t.Errorf("expected hello, got %s", result.String("APP_NAME"))
	}
}

func TestValidateMap_BooleanValues(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"1", true},
		{"yes", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"TRUE", true},
		{"YES", true},
	}
	for _, tc := range cases {
		v := envvalidator.New(
			envvalidator.Field{Key: "FEATURE_FLAG", Kind: envvalidator.KindBoolean},
		)
		result, err := v.ValidateMap(context.Background(), map[string]string{"FEATURE_FLAG": tc.input})
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
			continue
		}
		if result.Boolean("FEATURE_FLAG") != tc.expected {
			t.Errorf("input %q: expected %v, got %v", tc.input, tc.expected, result.Boolean("FEATURE_FLAG"))
		}
	}
}

func TestValidateMap_InvalidBoolean(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "FEATURE_FLAG", Kind: envvalidator.KindBoolean},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{"FEATURE_FLAG": "maybe"})
	if err == nil {
		t.Fatal("expected error for invalid boolean, got nil")
	}
}

func TestValidateMap_FloatKind(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "RATIO", Kind: envvalidator.KindFloat, Default: "0.5"},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{"RATIO": "3.14"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Float("RATIO") != 3.14 {
		t.Errorf("expected 3.14, got %f", result.Float("RATIO"))
	}
}

func TestValidateMap_URLKind(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "API_URL", Kind: envvalidator.KindURL, Required: true},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{"API_URL": "https://api.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String("API_URL") != "https://api.example.com" {
		t.Errorf("unexpected value: %s", result.String("API_URL"))
	}
}

func TestValidateMap_InvalidURL(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "API_URL", Kind: envvalidator.KindURL, Required: true},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{"API_URL": "not-a-url"})
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestValidateMap_DurationKind(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "TIMEOUT", Kind: envvalidator.KindDuration, Default: "30s"},
	)
	result, err := v.ValidateMap(context.Background(), map[string]string{"TIMEOUT": "2m"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := envvalidator.DurationResult(result, "TIMEOUT")
	if d != 2*time.Minute {
		t.Errorf("expected 2m, got %v", d)
	}
}

func TestValidateMap_InvalidDuration(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "TIMEOUT", Kind: envvalidator.KindDuration, Required: true},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{"TIMEOUT": "not-a-duration"})
	if err == nil {
		t.Fatal("expected error for invalid duration, got nil")
	}
}

func TestValidateMap_AllowedValues(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{
			Key:           "LOG_LEVEL",
			Kind:          envvalidator.KindString,
			Default:       "info",
			AllowedValues: []string{"debug", "info", "warn", "error"},
		},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{"LOG_LEVEL": "trace"})
	if err == nil {
		t.Fatal("expected error for disallowed value, got nil")
	}
	result, err := v.ValidateMap(context.Background(), map[string]string{"LOG_LEVEL": "debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String("LOG_LEVEL") != "debug" {
		t.Errorf("expected debug, got %s", result.String("LOG_LEVEL"))
	}
}

func TestValidateMap_MultipleErrors(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Required: true},
		envvalidator.Field{Key: "DATABASE_URL", Kind: envvalidator.KindURL, Required: true},
	)
	_, err := v.ValidateMap(context.Background(), map[string]string{})
	if err == nil {
		t.Fatal("expected errors, got nil")
	}
	verrs, ok := err.(envvalidator.ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors type, got %T", err)
	}
	if len(verrs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(verrs))
	}
}

func TestSchema_ReturnsAllFields(t *testing.T) {
	v := envvalidator.New(
		envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080", Description: "HTTP port"},
		envvalidator.Field{Key: "DATABASE_URL", Kind: envvalidator.KindURL, Required: true, Description: "Postgres URL"},
	)
	schema := v.Schema()
	if len(schema) != 2 {
		t.Fatalf("expected 2 schema entries, got %d", len(schema))
	}
	if schema[0].Key != "PORT" {
		t.Errorf("expected PORT, got %s", schema[0].Key)
	}
	if schema[1].Required != true {
		t.Error("expected DATABASE_URL to be required")
	}
}

func TestValidateMap_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	v := envvalidator.New(
		envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080"},
	)
	_, err := v.ValidateMap(ctx, map[string]string{})
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestValidationError_Error(t *testing.T) {
	e := &envvalidator.ValidationError{Key: "PORT", Reason: "required variable is missing or empty"}
	expected := `env-validator: field "PORT": required variable is missing or empty`
	if e.Error() != expected {
		t.Errorf("unexpected error string: %s", e.Error())
	}
}
