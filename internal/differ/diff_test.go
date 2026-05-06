package differ

import (
	"strings"
	"testing"
)

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}
	entries := Diff(base, target)

	found := findEntry(entries, "BAZ")
	if found == nil {
		t.Fatal("expected BAZ entry")
	}
	if found.Status != StatusAdded {
		t.Errorf("expected added, got %s", found.Status)
	}
	if found.NewValue != "qux" {
		t.Errorf("expected NewValue=qux, got %s", found.NewValue)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "val"}
	target := map[string]string{"FOO": "bar"}
	entries := Diff(base, target)

	found := findEntry(entries, "OLD")
	if found == nil {
		t.Fatal("expected OLD entry")
	}
	if found.Status != StatusRemoved {
		t.Errorf("expected removed, got %s", found.Status)
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"KEY": "old"}
	target := map[string]string{"KEY": "new"}
	entries := Diff(base, target)

	found := findEntry(entries, "KEY")
	if found == nil {
		t.Fatal("expected KEY entry")
	}
	if found.Status != StatusChanged {
		t.Errorf("expected changed, got %s", found.Status)
	}
	if found.OldValue != "old" || found.NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", found.OldValue, found.NewValue)
	}
}

func TestDiff_Same(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	target := map[string]string{"KEY": "val"}
	entries := Diff(base, target)

	found := findEntry(entries, "KEY")
	if found == nil || found.Status != StatusSame {
		t.Errorf("expected same status")
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "2"}
	target := map[string]string{"Z": "1", "A": "2"}
	entries := Diff(base, target)
	if entries[0].Key != "A" || entries[1].Key != "Z" {
		t.Errorf("expected sorted keys, got %s %s", entries[0].Key, entries[1].Key)
	}
}

func TestFormat_MaskSecrets(t *testing.T) {
	entries := []DiffEntry{
		{Key: "SECRET", Status: StatusAdded, NewValue: "mysecret"},
	}
	out := Format(entries, true)
	if strings.Contains(out, "mysecret") {
		t.Error("expected secret to be masked")
	}
	if !strings.Contains(out, "m******t") {
		t.Errorf("unexpected masked output: %s", out)
	}
}

func TestFormat_NoMask(t *testing.T) {
	entries := []DiffEntry{
		{Key: "KEY", Status: StatusChanged, OldValue: "old", NewValue: "new"},
	}
	out := Format(entries, false)
	if !strings.Contains(out, "old") || !strings.Contains(out, "new") {
		t.Errorf("expected plain values in output: %s", out)
	}
}

func findEntry(entries []DiffEntry, key string) *DiffEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
