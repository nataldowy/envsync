package syncer

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestSync_AddMissing(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=envsync\nNEW_KEY=hello\n")
	dst := writeTempEnv(t, "APP_NAME=envsync\n")

	res, err := Sync(src, dst, ModeAddMissing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", res.Added)
	}
	if len(res.Updated) != 0 {
		t.Errorf("expected no updates, got %v", res.Updated)
	}

	out := readFile(t, dst)
	if !contains(out, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY in output, got:\n%s", out)
	}
}

func TestSync_OverwriteChanged(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=newname\n")
	dst := writeTempEnv(t, "APP_NAME=oldname\n")

	res, err := Sync(src, dst, ModeOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Updated) != 1 || res.Updated[0] != "APP_NAME" {
		t.Errorf("expected Updated=[APP_NAME], got %v", res.Updated)
	}

	out := readFile(t, dst)
	if !contains(out, "APP_NAME=newname") {
		t.Errorf("expected APP_NAME=newname in output, got:\n%s", out)
	}
}

func TestSync_SkipChangedInAddMode(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=newname\n")
	dst := writeTempEnv(t, "APP_NAME=oldname\n")

	res, err := Sync(src, dst, ModeAddMissing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Skipped) != 1 || res.Skipped[0] != "APP_NAME" {
		t.Errorf("expected Skipped=[APP_NAME], got %v", res.Skipped)
	}
}

func TestSync_InvalidSource(t *testing.T) {
	_, err := Sync(filepath.Join(t.TempDir(), "missing.env"), writeTempEnv(t, ""), ModeAddMissing)
	if err == nil {
		t.Error("expected error for missing source file")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
