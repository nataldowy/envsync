package rotator

import (
	"strings"
	"testing"
)

func TestRotate_SensitiveKeysReplaced(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "old-pass",
		"API_KEY":     "old-key",
		"APP_NAME":    "myapp",
	}

	updated, res, err := Rotate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated["DB_PASSWORD"] == "old-pass" {
		t.Error("DB_PASSWORD should have been rotated")
	}
	if updated["API_KEY"] == "old-key" {
		t.Error("API_KEY should have been rotated")
	}
	if updated["APP_NAME"] != "myapp" {
		t.Error("APP_NAME should not have been rotated")
	}

	if _, ok := res.Rotated["DB_PASSWORD"]; !ok {
		t.Error("DB_PASSWORD missing from Rotated result")
	}
	if _, ok := res.Rotated["API_KEY"]; !ok {
		t.Error("API_KEY missing from Rotated result")
	}
}

func TestRotate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"SECRET_TOKEN": "original"}
	_, _, err := Rotate(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET_TOKEN"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestRotate_OnlyKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "pass",
		"API_SECRET":  "secret",
		"AUTH_TOKEN":  "token",
	}
	opts := DefaultOptions()
	opts.OnlyKeys = []string{"DB_PASSWORD"}

	updated, res, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated["DB_PASSWORD"] == "pass" {
		t.Error("DB_PASSWORD should have been rotated")
	}
	if updated["API_SECRET"] != "secret" {
		t.Error("API_SECRET should not have been rotated")
	}
	if len(res.Rotated) != 1 {
		t.Errorf("expected 1 rotated key, got %d", len(res.Rotated))
	}
}

func TestRotate_Prefix(t *testing.T) {
	env := map[string]string{"APP_SECRET": "old"}
	opts := DefaultOptions()
	opts.Prefix = "rot_"

	updated, _, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(updated["APP_SECRET"], "rot_") {
		t.Errorf("expected value to start with prefix, got %q", updated["APP_SECRET"])
	}
}

func TestRotate_NonDeterministic(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "same"}

	updated1, _, _ := Rotate(env, DefaultOptions())
	updated2, _, _ := Rotate(env, DefaultOptions())

	if updated1["DB_PASSWORD"] == updated2["DB_PASSWORD"] {
		t.Error("two rotations produced the same value (collision is astronomically unlikely)")
	}
}

func TestRotate_DefaultByteLength(t *testing.T) {
	env := map[string]string{"API_KEY": "x"}
	opts := Options{ByteLength: 0} // should default to 32

	updated, _, err := Rotate(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// hex of 32 bytes = 64 chars
	if len(updated["API_KEY"]) != 64 {
		t.Errorf("expected 64-char hex value, got len=%d", len(updated["API_KEY"]))
	}
}
