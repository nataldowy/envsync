package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/differ"
	"github.com/yourorg/envsync/internal/reporter"
)

func TestReporter_PlainFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatPlain)

	entries := []differ.Entry{
		{Key: "DB_HOST", Status: differ.Added, RightVal: "localhost"},
		{Key: "API_KEY", Status: differ.Removed, LeftVal: "secret"},
		{Key: "PORT", Status: differ.Changed, LeftVal: "3000", RightVal: "8080"},
		{Key: "APP_NAME", Status: differ.Same, LeftVal: "myapp", RightVal: "myapp"},
	}

	err := r.Write(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "+") {
		t.Error("expected added marker '+' in plain output")
	}
	if !strings.Contains(out, "-") {
		t.Error("expected removed marker '-' in plain output")
	}
	if !strings.Contains(out, "~") {
		t.Error("expected changed marker '~' in plain output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST in output")
	}
}

func TestReporter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)

	entries := []differ.Entry{
		{Key: "DB_HOST", Status: differ.Added, RightVal: "localhost"},
		{Key: "API_KEY", Status: differ.Removed, LeftVal: "secret"},
	}

	err := r.Write(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"key"`) {
		t.Error("expected JSON key field")
	}
	if !strings.Contains(out, `"status"`) {
		t.Error("expected JSON status field")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in JSON output")
	}
}

func TestReporter_Summary(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatPlain)

	entries := []differ.Entry{
		{Key: "A", Status: differ.Added, RightVal: "1"},
		{Key: "B", Status: differ.Added, RightVal: "2"},
		{Key: "C", Status: differ.Removed, LeftVal: "3"},
		{Key: "D", Status: differ.Same, LeftVal: "4", RightVal: "4"},
	}

	_ = r.Write(entries)
	s := r.Summary()

	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Same != 1 {
		t.Errorf("expected 1 same, got %d", s.Same)
	}
}

func TestReporter_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatPlain)

	err := r.Write([]differ.Entry{})
	if err != nil {
		t.Fatalf("unexpected error on empty entries: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "no differences") {
		t.Errorf("expected 'no differences' message, got: %q", out)
	}
}
