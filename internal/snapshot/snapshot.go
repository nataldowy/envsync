// Package snapshot provides functionality to capture and compare
// .env file states at different points in time.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured state of an .env file.
type Snapshot struct {
	File      string            `json:"file"`
	CapturedAt time.Time        `json:"captured_at"`
	Entries   map[string]string `json:"entries"`
}

// Capture reads the given env map and returns a Snapshot.
func Capture(file string, entries map[string]string) *Snapshot {
	copy := make(map[string]string, len(entries))
	for k, v := range entries {
		copy[k] = v
	}
	return &Snapshot{
		File:       file,
		CapturedAt: time.Now().UTC(),
		Entries:    copy,
	}
}

// Save writes the snapshot as JSON to the given path.
func Save(s *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a snapshot from the given JSON file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode %q: %w", path, err)
	}
	return &s, nil
}

// Delta describes the difference between two snapshots.
type Delta struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// Compare returns the Delta between an older and a newer snapshot.
func Compare(old, new *Snapshot) Delta {
	d := Delta{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, v := range new.Entries {
		if ov, ok := old.Entries[k]; !ok {
			d.Added[k] = v
		} else if ov != v {
			d.Changed[k] = [2]string{ov, v}
		}
	}
	for k, v := range old.Entries {
		if _, ok := new.Entries[k]; !ok {
			d.Removed[k] = v
		}
	}
	return d
}
