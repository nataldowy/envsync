package linter_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/linter"
)

func makeIndex(keys []string) map[string]int {
	idx := make(map[string]int, len(keys))
	for i, k := range keys {
		idx[k] = i + 1
	}
	return idx
}

func TestLint_UpperCase(t *testing.T) {
	entries := map[string]string{"my_key": "value", "GOOD_KEY": "value"}
	idx := makeIndex([]string{"my_key", "GOOD_KEY"})
	opts := linter.DefaultOptions()

	issues := linter.Lint(entries, idx, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "my_key" {
		t.Errorf("expected issue on my_key, got %s", issues[0].Key)
	}
	if issues[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestLint_LeadingDigit(t *testing.T) {
	entries := map[string]string{"1BAD": "val"}
	idx := makeIndex([]string{"1BAD"})
	opts := linter.DefaultOptions()

	issues := linter.Lint(entries, idx, opts)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for leading-digit key")
	}
	var found bool
	for _, iss := range issues {
		if iss.Severity == linter.SeverityError && iss.Key == "1BAD" {
			found = true
		}
	}
	if !found {
		t.Error("expected error-severity issue for key starting with digit")
	}
}

func TestLint_EmptyValueForbidden(t *testing.T) {
	entries := map[string]string{"EMPTY_KEY": ""}
	idx := makeIndex([]string{"EMPTY_KEY"})
	opts := linter.DefaultOptions()
	opts.AllowEmptyValues = false

	issues := linter.Lint(entries, idx, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity for empty value")
	}
}

func TestLint_EmptyValueAllowed(t *testing.T) {
	entries := map[string]string{"EMPTY_KEY": ""}
	idx := makeIndex([]string{"EMPTY_KEY"})
	opts := linter.DefaultOptions()
	opts.AllowEmptyValues = true
	opts.RequireUpperCase = true // EMPTY_KEY is already upper

	issues := linter.Lint(entries, idx, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_MaxKeyLength(t *testing.T) {
	longKey := "ABCDEFGHIJKLMNOPQRSTUVWXYZ_ABCDEFGHIJKLMNOPQRSTUVWXYZ_TOOLONG"
	entries := map[string]string{longKey: "val"}
	idx := makeIndex([]string{longKey})
	opts := linter.DefaultOptions()
	opts.MaxKeyLength = 20

	issues := linter.Lint(entries, idx, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue for long key, got %d", len(issues))
	}
}

func TestLint_NoIssues(t *testing.T) {
	entries := map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"}
	idx := makeIndex([]string{"DB_HOST", "APP_PORT"})
	opts := linter.DefaultOptions()

	issues := linter.Lint(entries, idx, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}
