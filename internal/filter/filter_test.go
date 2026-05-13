package filter_test

import (
	"testing"

	"github.com/user/envsync/internal/filter"
)

var base = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envsync",
	"APP_VERSION": "1.0.0",
	"LOG_LEVEL":   "info",
	"SECRET_KEY":  "abc123",
}

func TestFilter_Prefix(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	for k := range out {
		if k != "DB_HOST" && k != "DB_PORT" {
			t.Errorf("unexpected key %q", k)
		}
	}
}

func TestFilter_Suffix(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Suffix: "_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["SECRET_KEY"]; !ok {
		t.Error("expected SECRET_KEY in result")
	}
}

func TestFilter_Pattern(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{Pattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.Filter(base, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestFilter_AllowList(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{AllowList: []string{"LOG_LEVEL", "APP_NAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_DenyList(t *testing.T) {
	out, err := filter.Filter(base, filter.Options{DenyList: []string{"SECRET_KEY", "DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("SECRET_KEY should have been removed")
	}
	if _, ok := out["DB_PORT"]; ok {
		t.Error("DB_PORT should have been removed")
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(out))
	}
}

func TestFilter_NoOptions(t *testing.T) {
	out, err := filter.Filter(base, filter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(base) {
		t.Fatalf("expected %d keys, got %d", len(base), len(out))
	}
}
