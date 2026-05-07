package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envsync/internal/snapshot"
)

func TestCapture_CopiesEntries(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.Capture(".env", original)
	original["FOO"] = "mutated"
	if s.Entries["FOO"] != "bar" {
		t.Errorf("expected snapshot to be isolated from source map")
	}
	if s.File != ".env" {
		t.Errorf("expected file field to be set")
	}
}

func TestSaveLoad_Roundtrip(t *testing.T) {
	entries := map[string]string{"KEY": "value", "SECRET": "s3cr3t"}
	s := snapshot.Capture(".env", entries)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.File != s.File {
		t.Errorf("file mismatch: got %q want %q", loaded.File, s.File)
	}
	for k, v := range entries {
		if loaded.Entries[k] != v {
			t.Errorf("key %q: got %q want %q", k, loaded.Entries[k], v)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCompare_AddedRemovedChanged(t *testing.T) {
	old := snapshot.Capture(".env", map[string]string{
		"KEEP":    "same",
		"CHANGE":  "old",
		"REMOVED": "gone",
	})
	new := snapshot.Capture(".env", map[string]string{
		"KEEP":   "same",
		"CHANGE": "new",
		"ADDED":  "fresh",
	})

	d := snapshot.Compare(old, new)

	if _, ok := d.Added["ADDED"]; !ok {
		t.Error("expected ADDED in Added")
	}
	if _, ok := d.Removed["REMOVED"]; !ok {
		t.Error("expected REMOVED in Removed")
	}
	pair, ok := d.Changed["CHANGE"]
	if !ok {
		t.Fatal("expected CHANGE in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected changed pair: %v", pair)
	}
	if _, ok := d.Changed["KEEP"]; ok {
		t.Error("KEEP should not appear in Changed")
	}
}

func TestCompare_EmptySnapshots(t *testing.T) {
	old := snapshot.Capture(".env", map[string]string{})
	new := snapshot.Capture(".env", map[string]string{})
	d := snapshot.Compare(old, new)
	if len(d.Added)+len(d.Removed)+len(d.Changed) != 0 {
		t.Error("expected empty delta for identical empty snapshots")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	s := snapshot.Capture(".env", map[string]string{"A": "1"})
	err := snapshot.Save(s, filepath.Join(t.TempDir(), "no", "such", "dir", "snap.json"))
	if err == nil {
		t.Fatal("expected error for invalid save path")
	}
	_ = os.RemoveAll("no")
}
