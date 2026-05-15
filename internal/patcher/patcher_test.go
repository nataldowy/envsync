package patcher_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/patcher"
)

func base() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
	}
}

func TestApply_Set_NewKey(t *testing.T) {
	r, err := patcher.Apply(base(), []patcher.Op{{Kind: patcher.OpSet, Key: "NEW_VAR", Value: "hello"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.Env["NEW_VAR"] != "hello" {
		t.Errorf("expected NEW_VAR=hello, got %q", r.Env["NEW_VAR"])
	}
	if len(r.Changes) != 1 || r.Changes[0].Op != patcher.OpSet {
		t.Errorf("unexpected changes: %v", r.Changes)
	}
}

func TestApply_Set_OverwriteExisting(t *testing.T) {
	r, err := patcher.Apply(base(), []patcher.Op{{Kind: patcher.OpSet, Key: "DB_HOST", Value: "prod-db"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.Env["DB_HOST"] != "prod-db" {
		t.Errorf("expected prod-db, got %q", r.Env["DB_HOST"])
	}
	if r.Changes[0].OldValue != "localhost" {
		t.Errorf("expected OldValue=localhost, got %q", r.Changes[0].OldValue)
	}
}

func TestApply_Delete(t *testing.T) {
	r, err := patcher.Apply(base(), []patcher.Op{{Kind: patcher.OpDelete, Key: "API_KEY"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Env["API_KEY"]; ok {
		t.Error("API_KEY should have been deleted")
	}
}

func TestApply_Delete_MissingKey_Error(t *testing.T) {
	_, err := patcher.Apply(base(), []patcher.Op{{Kind: patcher.OpDelete, Key: "MISSING"}})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestApply_Rename(t *testing.T) {
	r, err := patcher.Apply(base(), []patcher.Op{{Kind: patcher.OpRename, Key: "DB_HOST", NewKey: "DATABASE_HOST"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := r.Env["DB_HOST"]; ok {
		t.Error("old key DB_HOST should not exist")
	}
	if r.Env["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", r.Env["DATABASE_HOST"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	src := base()
	patcher.Apply(src, []patcher.Op{ //nolint:errcheck
		{Kind: patcher.OpSet, Key: "DB_HOST", Value: "changed"},
	})
	if src["DB_HOST"] != "localhost" {
		t.Error("source map was mutated")
	}
}

func TestApply_UnknownOp_Error(t *testing.T) {
	_, err := patcher.Apply(base(), []patcher.Op{{Kind: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestFormatResult_ContainsKeys(t *testing.T) {
	r, _ := patcher.Apply(base(), []patcher.Op{
		{Kind: patcher.OpSet, Key: "DB_HOST", Value: "prod"},
		{Kind: patcher.OpDelete, Key: "API_KEY"},
	})
	out := patcher.FormatResult(r, nil)
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
}

func TestSummary_Counts(t *testing.T) {
	r, _ := patcher.Apply(base(), []patcher.Op{
		{Kind: patcher.OpSet, Key: "NEW", Value: "v"},
		{Kind: patcher.OpDelete, Key: "API_KEY"},
		{Kind: patcher.OpRename, Key: "DB_HOST", NewKey: "DATABASE_HOST"},
	})
	s := patcher.Summary(r)
	if !strings.Contains(s, "1 set") || !strings.Contains(s, "1 deleted") || !strings.Contains(s, "1 renamed") {
		t.Errorf("unexpected summary: %q", s)
	}
}
