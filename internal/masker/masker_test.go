package masker

import (
	"testing"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := New(nil)

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AWS_SECRET_ACCESS_KEY",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"MY_CREDENTIAL",
	}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}

	insensitive := []string{
		"APP_ENV",
		"PORT",
		"DEBUG",
		"LOG_LEVEL",
	}
	for _, key := range insensitive {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	m := New([]string{"MAGIC", "HIDDEN"})

	if !m.IsSensitive("MAGIC_VALUE") {
		t.Error("expected MAGIC_VALUE to be sensitive")
	}
	if !m.IsSensitive("MY_HIDDEN_KEY") {
		t.Error("expected MY_HIDDEN_KEY to be sensitive")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom patterns")
	}
}

func TestMask_SensitiveValue(t *testing.T) {
	m := New(nil)
	got := m.Mask("DB_PASSWORD", "supersecret")
	if got != MaskedValue {
		t.Errorf("expected %q, got %q", MaskedValue, got)
	}
}

func TestMask_PlainValue(t *testing.T) {
	m := New(nil)
	got := m.Mask("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected %q, got %q", "production", got)
	}
}

func TestMaskMap(t *testing.T) {
	m := New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "key-abc-123",
		"PORT":        "8080",
	}

	masked := m.MaskMap(env)

	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unmasked, got %q", masked["APP_ENV"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should be unmasked, got %q", masked["PORT"])
	}
	if masked["DB_PASSWORD"] != MaskedValue {
		t.Errorf("DB_PASSWORD should be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_KEY"] != MaskedValue {
		t.Errorf("API_KEY should be masked, got %q", masked["API_KEY"])
	}
	// original map must be unchanged
	if env["DB_PASSWORD"] != "s3cr3t" {
		t.Error("original map was mutated")
	}
}
