package redactor_test

import (
	"regexp"
	"testing"

	"github.com/user/envsync/internal/redactor"
)

func defaultRedactor() *redactor.Redactor {
	return redactor.New(redactor.Options{})
}

func TestIsSensitive_MatchesDefaults(t *testing.T) {
	r := defaultRedactor()
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY", "SECRET"}
	for _, k := range sensitive {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_PlainKey(t *testing.T) {
	r := defaultRedactor()
	if r.IsSensitive("APP_ENV") {
		t.Error("APP_ENV should not be sensitive")
	}
}

func TestIsSensitive_ExtraPattern(t *testing.T) {
	r := redactor.New(redactor.Options{
		ExtraPattern: regexp.MustCompile(`(?i)pin`),
	})
	if !r.IsSensitive("USER_PIN") {
		t.Error("USER_PIN should match extra pattern")
	}
}

func TestRedactMap_SensitiveReplaced(t *testing.T) {
	r := defaultRedactor()
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_ENV":     "production",
	}
	out := r.RedactMap(env)
	if out["DB_PASSWORD"] != redactor.DefaultPlaceholder {
		t.Errorf("expected placeholder, got %q", out["DB_PASSWORD"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("plain value should be unchanged, got %q", out["APP_ENV"])
	}
}

func TestRedactMap_CustomPlaceholder(t *testing.T) {
	r := redactor.New(redactor.Options{Placeholder: "***"})
	out := r.RedactMap(map[string]string{"SECRET_KEY": "abc"})
	if out["SECRET_KEY"] != "***" {
		t.Errorf("expected ***, got %q", out["SECRET_KEY"])
	}
}

func TestRedactLine_SensitiveLine(t *testing.T) {
	r := defaultRedactor()
	got := r.RedactLine("API_KEY=supersecret")
	if got != "API_KEY="+redactor.DefaultPlaceholder {
		t.Errorf("unexpected redacted line: %q", got)
	}
}

func TestRedactLine_PlainLine(t *testing.T) {
	r := defaultRedactor()
	got := r.RedactLine("APP_ENV=production")
	if got != "APP_ENV=production" {
		t.Errorf("plain line should be unchanged, got %q", got)
	}
}

func TestRedactLine_CommentPassthrough(t *testing.T) {
	r := defaultRedactor()
	got := r.RedactLine("# this is a comment")
	if got != "# this is a comment" {
		t.Errorf("comment should pass through unchanged, got %q", got)
	}
}

func TestRedactLines_MultiLine(t *testing.T) {
	r := defaultRedactor()
	src := "APP_ENV=dev\nDB_PASSWORD=secret\n# comment\nPORT=8080"
	out := r.RedactLines(src)
	expected := "APP_ENV=dev\nDB_PASSWORD=" + redactor.DefaultPlaceholder + "\n# comment\nPORT=8080"
	if out != expected {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", out, expected)
	}
}
