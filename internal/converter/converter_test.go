package converter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/converter"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"DB_PORT":  "5432",
}

func TestConvert_Export(t *testing.T) {
	out, err := converter.Convert(sampleEnv, converter.FormatExport, converter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export statement, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_PORT=5432") {
		t.Errorf("expected DB_PORT export, got:\n%s", out)
	}
}

func TestConvert_Docker(t *testing.T) {
	out, err := converter.Convert(sampleEnv, converter.FormatDocker, converter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--env DB_HOST=localhost") {
		t.Errorf("expected docker flag, got:\n%s", out)
	}
}

func TestConvert_Inline(t *testing.T) {
	out, err := converter.Convert(sampleEnv, converter.FormatInline, converter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// inline is space-separated; check no newlines mid-content
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected single line, got %d lines", len(lines))
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected inline pair, got: %s", out)
	}
}

func TestConvert_Makefile(t *testing.T) {
	out, err := converter.Convert(sampleEnv, converter.FormatMakefile, converter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_PORT := 5432") {
		t.Errorf("expected makefile assignment, got:\n%s", out)
	}
}

func TestConvert_QuoteValues(t *testing.T) {
	opts := converter.DefaultOptions()
	opts.QuoteValues = true
	out, err := converter.Convert(map[string]string{"KEY": "val"}, converter.FormatExport, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `export KEY="val"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestConvert_SortedOutput(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	out, err := converter.Convert(env, converter.FormatExport, converter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "export A_KEY") {
		t.Errorf("expected A_KEY first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "export Z_KEY") {
		t.Errorf("expected Z_KEY last, got: %s", lines[2])
	}
}

func TestConvert_UnsupportedFormat(t *testing.T) {
	_, err := converter.Convert(sampleEnv, converter.Format("toml"), converter.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
