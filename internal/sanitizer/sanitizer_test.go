package sanitizer_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/sanitizer"
)

func TestSanitize_UpperCaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "Api_Key": "abc"}
	out, err := sanitizer.Sanitize(src, sanitizer.Options{UpperCaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"DB_HOST", "API_KEY"} {
		if _, ok := out[k]; !ok {
			t.Errorf("expected key %q in output", k)
		}
	}
}

func TestSanitize_TrimValues(t *testing.T) {
	src := map[string]string{"KEY": "  hello world  "}
	out, err := sanitizer.Sanitize(src, sanitizer.Options{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["KEY"]; got != "hello world" {
		t.Errorf("expected trimmed value, got %q", got)
	}
}

func TestSanitize_StripInvalidKeyChars(t *testing.T) {
	src := map[string]string{"my-key": "v1", "another.key": "v2"}
	out, err := sanitizer.Sanitize(src, sanitizer.Options{StripInvalidKeyChars: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["mykey"]; !ok {
		t.Error("expected key 'mykey' after stripping hyphen")
	}
	if _, ok := out["anotherkey"]; !ok {
		t.Error("expected key 'anotherkey' after stripping dot")
	}
}

func TestSanitize_DropEmptyKeyAfterStrip(t *testing.T) {
	src := map[string]string{"---": "value"}
	out, err := sanitizer.Sanitize(src, sanitizer.Options{StripInvalidKeyChars: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestSanitize_NoMutationOfInput(t *testing.T) {
	src := map[string]string{"key": "  val  "}
	_, _ = sanitizer.Sanitize(src, sanitizer.DefaultOptions())
	if src["key"] != "  val  " {
		t.Error("Sanitize must not mutate the input map")
	}
}

func TestSanitize_DefaultOptions(t *testing.T) {
	src := map[string]string{"db-host": "  127.0.0.1  "}
	out, err := sanitizer.Sanitize(src, sanitizer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// hyphen stripped, key upper-cased, value trimmed
	if v, ok := out["DBHOST"]; !ok || v != "127.0.0.1" {
		t.Errorf("unexpected result: %v", out)
	}
}
