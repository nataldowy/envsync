package deduplicator

import (
	"testing"
)

func TestDeduplicate_NoDuplicates(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux"}
	res, err := Deduplicate(lines, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
	if res.Env["FOO"] != "bar" || res.Env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestDeduplicate_StrategyLast(t *testing.T) {
	lines := []string{"KEY=first", "OTHER=x", "KEY=second"}
	res, err := Deduplicate(lines, Options{Strategy: StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", res.Env["KEY"])
	}
	if len(res.Duplicates["KEY"]) != 2 {
		t.Errorf("expected 2 values recorded, got %v", res.Duplicates["KEY"])
	}
}

func TestDeduplicate_StrategyFirst(t *testing.T) {
	lines := []string{"KEY=first", "KEY=second", "KEY=third"}
	res, err := Deduplicate(lines, Options{Strategy: StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
}

func TestDeduplicate_StrategyError(t *testing.T) {
	lines := []string{"KEY=a", "KEY=b"}
	_, err := Deduplicate(lines, Options{Strategy: StrategyError})
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestDeduplicate_StrategyError_NoDuplicates(t *testing.T) {
	lines := []string{"ALPHA=1", "BETA=2"}
	res, err := Deduplicate(lines, Options{Strategy: StrategyError})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestDeduplicate_IgnoresCommentsAndBlanks(t *testing.T) {
	lines := []string{"# comment", "", "FOO=bar"}
	res, err := Deduplicate(lines, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 1 || res.Env["FOO"] != "bar" {
		t.Errorf("unexpected env: %v", res.Env)
	}
}

func TestDeduplicate_DefaultStrategyIsLast(t *testing.T) {
	lines := []string{"X=1", "X=2"}
	res, err := Deduplicate(lines, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["X"] != "2" {
		t.Errorf("expected '2' with default strategy, got %q", res.Env["X"])
	}
}
