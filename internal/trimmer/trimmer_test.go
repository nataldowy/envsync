package trimmer

import (
	"testing"
)

func TestTrim_TrimSpaceDefault(t *testing.T) {
	src := map[string]string{
		"  KEY  ": "  value  ",
	}
	out, err := Trim(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["KEY"]; !ok || v != "value" {
		t.Errorf("expected KEY=value, got %q=%q", "KEY", v)
	}
}

func TestTrim_RemoveEmpty(t *testing.T) {
	src := map[string]string{
		"PRESENT": "hello",
		"EMPTY":   "",
		"BLANK":   "   ",
	}
	opts := DefaultOptions()
	opts.RemoveEmpty = true
	out, err := Trim(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["EMPTY"]; ok {
		t.Error("EMPTY should have been removed")
	}
	if _, ok := out["BLANK"]; ok {
		t.Error("BLANK should have been removed")
	}
	if out["PRESENT"] != "hello" {
		t.Errorf("PRESENT should be kept, got %q", out["PRESENT"])
	}
}

func TestTrim_KeepEmptyWhenFlagOff(t *testing.T) {
	src := map[string]string{"EMPTY": ""}
	opts := DefaultOptions()
	opts.RemoveEmpty = false
	out, _ := Trim(src, opts)
	if _, ok := out["EMPTY"]; !ok {
		t.Error("EMPTY should be retained when RemoveEmpty=false")
	}
}

func TestTrim_StripPrefix(t *testing.T) {
	src := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DEBUG":    "true",
	}
	opts := DefaultOptions()
	opts.StripPrefixes = []string{"APP_"}
	out, err := Trim(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_")
	}
	if _, ok := out["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_")
	}
	if _, ok := out["DEBUG"]; !ok {
		t.Error("DEBUG has no prefix and should be kept as-is")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("APP_HOST should not appear after stripping")
	}
}

func TestTrim_EmptyKeyDropped(t *testing.T) {
	src := map[string]string{
		"APP_": "orphan",
	}
	opts := DefaultOptions()
	opts.StripPrefixes = []string{"APP_"}
	out, err := Trim(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	opts := DefaultOptions()
	opts.StripPrefixes = []string{"KEY"}
	Trim(src, opts)
	if _, ok := src["KEY"]; !ok {
		t.Error("Trim must not mutate the source map")
	}
}
