// Package auditor records and retrieves a history of sync and diff operations
// performed by envsync, enabling traceability across environments.
package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// EventType classifies the kind of operation that was audited.
type EventType string

const (
	EventDiff EventType = "diff"
	EventSync EventType = "sync"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	Changes   int       `json:"changes"`
	Message   string    `json:"message,omitempty"`
}

// Auditor writes audit entries to a JSON-lines log file.
type Auditor struct {
	path string
}

// New returns an Auditor that appends entries to the file at path.
func New(path string) *Auditor {
	return &Auditor{path: path}
}

// Record appends a new audit entry to the log file.
func (a *Auditor) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("auditor: open log file: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("auditor: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll parses and returns all entries stored in the audit log.
func (a *Auditor) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(a.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("auditor: read log file: %w", err)
	}

	var entries []Entry
	decoder := json.NewDecoder(
		newBytesReader(data),
	)
	for decoder.More() {
		var e Entry
		if err := decoder.Decode(&e); err != nil {
			return nil, fmt.Errorf("auditor: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
