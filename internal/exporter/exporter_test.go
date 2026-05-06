package exporter

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestExport_DotEnv(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	tmp := t.TempDir() + "/out.env"
	if err := Export(env, tmp, FormatDotEnv, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	content := string(data)
	if !strings.Contains(content, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", content)
	}
	if !strings.Contains(content, "PORT=8080") {
		t.Errorf("expected PORT=8080 in output, got:\n%s", content)
	}
}

func TestExport_JSON(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	tmp := t.TempDir() + "/out.json"
	if err := Export(env, tmp, FormatJSON, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", result["KEY"])
	}
}

func TestExport_YAML(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	tmp := t.TempDir() + "/out.yaml"
	if err := Export(env, tmp, FormatYAML, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "DB_HOST") {
		t.Errorf("expected DB_HOST in YAML output")
	}
}

func TestExport_MaskKeys(t *testing.T) {
	env := map[string]string{"API_KEY": "secret123", "PORT": "9000"}
	tmp := t.TempDir() + "/out.env"
	if err := Export(env, tmp, FormatDotEnv, []string{"API_KEY"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if strings.Contains(string(data), "secret123") {
		t.Errorf("expected secret to be masked, got: %s", string(data))
	}
	if !strings.Contains(string(data), "***") {
		t.Errorf("expected *** mask in output")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	env := map[string]string{"X": "1"}
	err := Export(env, "/tmp/x", Format("toml"), nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_DotEnvQuotesSpaces(t *testing.T) {
	env := map[string]string{"GREETING": "hello world"}
	tmp := t.TempDir() + "/out.env"
	if err := Export(env, tmp, FormatDotEnv, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), `"hello world"`) {
		t.Errorf("expected quoted value with space, got: %s", string(data))
	}
}
