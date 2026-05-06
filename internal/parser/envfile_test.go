package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "APP_ENV" || ef.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", ef.Entries[0])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"
SECRET='mysecret'
`)
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("expected unquoted DB_URL, got %q", m["DB_URL"])
	}
	if m["SECRET"] != "mysecret" {
		t.Errorf("expected unquoted SECRET, got %q", m["SECRET"])
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	content := "# database config\nDB_HOST=localhost\n\nPORT=5432\n"
	path := writeTempEnv(t, content)
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Comment != "database config" {
		t.Errorf("expected comment on DB_HOST, got %q", ef.Entries[0].Comment)
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "TIMEOUT=30 # seconds\n")
	ef, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Value != "30" {
		t.Errorf("expected value '30', got %q", ef.Entries[0].Value)
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestToMap(t *testing.T) {
	path := writeTempEnv(t, "A=1\nB=2\nA=3\n")
	ef, _ := Parse(path)
	m := ef.ToMap()
	// last value wins
	if m["A"] != "3" {
		t.Errorf("expected A=3 (last wins), got %q", m["A"])
	}
	if m["B"] != "2" {
		t.Errorf("expected B=2, got %q", m["B"])
	}
}
