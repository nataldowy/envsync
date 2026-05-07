package resolver

import (
	"testing"
)

func TestResolve_NoInterpolation(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("values changed unexpectedly: %v", out)
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{"BASE": "/usr/local", "BIN": "${BASE}/bin"}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BIN"] != "/usr/local/bin" {
		t.Errorf("got %q, want %q", out["BIN"], "/usr/local/bin")
	}
}

func TestResolve_NoBraceStyle(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "URL": "http://$HOST:8080"}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Errorf("got %q, want %q", out["URL"], "http://localhost:8080")
	}
}

func TestResolve_ChainedVars(t *testing.T) {
	env := map[string]string{"A": "hello", "B": "${A}_world", "C": "${B}!"}
	out, err := Resolve(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello_world!" {
		t.Errorf("got %q, want %q", out["C"], "hello_world!")
	}
}

func TestResolve_UndefinedVar_Error(t *testing.T) {
	env := map[string]string{"FOO": "${MISSING}"}
	_, err := Resolve(env, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for undefined variable, got nil")
	}
}

func TestResolve_UndefinedVar_AllowMissing(t *testing.T) {
	opts := DefaultOptions()
	opts.AllowMissing = true
	env := map[string]string{"FOO": "${MISSING}_suffix"}
	out, err := Resolve(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "_suffix" {
		t.Errorf("got %q, want %q", out["FOO"], "_suffix")
	}
}

func TestResolve_MaxDepthExceeded(t *testing.T) {
	opts := DefaultOptions()
	opts.MaxDepth = 2
	// A -> B -> C -> D exceeds depth 2
	env := map[string]string{"A": "${B}", "B": "${C}", "C": "${D}", "D": "end"}
	_, err := Resolve(env, opts)
	if err == nil {
		t.Fatal("expected max depth error, got nil")
	}
}
