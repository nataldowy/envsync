package auditor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestRecord_CreatesFile(t *testing.T) {
	path := tempLog(t)
	a := New(path)

	if err := a.Record(Entry{Event: EventDiff, Source: ".env", Target: ".env.staging", Changes: 2}); err != nil {
		t.Fatalf("Record: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("log file not created: %v", err)
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	path := tempLog(t)
	a := New(path)

	events := []Entry{
		{Event: EventDiff, Source: ".env", Target: ".env.prod", Changes: 1},
		{Event: EventSync, Source: ".env.prod", Target: ".env.local", Changes: 4},
	}
	for _, e := range events {
		if err := a.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	got, err := a.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Event != EventDiff {
		t.Errorf("entry[0].Event = %q, want %q", got[0].Event, EventDiff)
	}
	if got[1].Changes != 4 {
		t.Errorf("entry[1].Changes = %d, want 4", got[1].Changes)
	}
}

func TestRecord_TimestampAutoSet(t *testing.T) {
	path := tempLog(t)
	a := New(path)

	before := time.Now().UTC()
	if err := a.Record(Entry{Event: EventSync, Source: "a", Target: "b"}); err != nil {
		t.Fatalf("Record: %v", err)
	}
	after := time.Now().UTC()

	entries, _ := a.ReadAll()
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}

func TestReadAll_EmptyFile(t *testing.T) {
	path := tempLog(t)
	a := New(path)

	entries, err := a.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_MessageField(t *testing.T) {
	path := tempLog(t)
	a := New(path)

	const msg = "manual override applied"
	_ = a.Record(Entry{Event: EventSync, Source: "x", Target: "y", Message: msg})

	entries, _ := a.ReadAll()
	if entries[0].Message != msg {
		t.Errorf("Message = %q, want %q", entries[0].Message, msg)
	}
}
