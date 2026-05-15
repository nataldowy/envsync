package normalizer

import (
	"testing"
)

func TestNormalize_UpperCaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "api_key": "secret"}
	out := Normalize(env, Options{UpperCaseKeys: true})
	for _, want := range []string{"DB_HOST", "API_KEY"} {
		if _, ok := out[want]; !ok {
			t.Errorf("expected key %q in output", want)
		}
	}
	if _, ok := out["db_host"]; ok {
		t.Error("lower-case key should have been removed")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"HOST": "  localhost  ", "PORT": "\t5432\n"}
	out := Normalize(env, Options{TrimValues: true})
	if got := out["HOST"]; got != "localhost" {
		t.Errorf("HOST: got %q, want %q", got, "localhost")
	}
	if got := out["PORT"]; got != "5432" {
		t.Errorf("PORT: got %q, want %q", got, "5432")
	}
}

func TestNormalize_AddPrefix(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	out := Normalize(env, Options{AddPrefix: "APP_"})
	if _, ok := out["APP_HOST"]; !ok {
		t.Error("expected key APP_HOST")
	}
	if _, ok := out["HOST"]; ok {
		t.Error("original key HOST should not exist")
	}
}

func TestNormalize_StripPrefix(t *testing.T) {
	env := map[string]string{"PROD_HOST": "example.com", "PROD_PORT": "443"}
	out := Normalize(env, Options{StripPrefix: "PROD_"})
	for _, want := range []string{"HOST", "PORT"} {
		if _, ok := out[want]; !ok {
			t.Errorf("expected key %q after strip", want)
		}
	}
}

func TestNormalize_StripAndAddPrefix(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	out := Normalize(env, Options{StripPrefix: "OLD_", AddPrefix: "NEW_"})
	if _, ok := out["NEW_KEY"]; !ok {
		t.Error("expected NEW_KEY")
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"host": "  val  "}
	_ = Normalize(env, DefaultOptions())
	if env["host"] != "  val  " {
		t.Error("input map was mutated")
	}
}

func TestNormalize_EmptyKeyDroppedAfterStrip(t *testing.T) {
	// Stripping a prefix that equals the whole key leaves an empty string.
	env := map[string]string{"PREFIX": "value"}
	out := Normalize(env, Options{StripPrefix: "PREFIX"})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %v", out)
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if !opts.UpperCaseKeys {
		t.Error("DefaultOptions should enable UpperCaseKeys")
	}
	if !opts.TrimValues {
		t.Error("DefaultOptions should enable TrimValues")
	}
}
