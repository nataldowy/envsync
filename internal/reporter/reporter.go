// Package reporter provides functionality for generating structured reports
// of env file diff and sync operations, supporting multiple output formats.
package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Format represents the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// OperationType describes the kind of operation recorded in a report entry.
type OperationType string

const (
	OpAdded   OperationType = "added"
	OpRemoved OperationType = "removed"
	OpChanged OperationType = "changed"
	OpSame    OperationType = "same"
	OpSynced  OperationType = "synced"
	OpSkipped OperationType = "skipped"
)

// Entry represents a single key-level event in a report.
type Entry struct {
	Key       string        `json:"key"       yaml:"key"`
	Operation OperationType `json:"operation" yaml:"operation"`
	OldValue  string        `json:"old_value,omitempty" yaml:"old_value,omitempty"`
	NewValue  string        `json:"new_value,omitempty" yaml:"new_value,omitempty"`
	Masked    bool          `json:"masked"    yaml:"masked"`
}

// Report holds metadata and all entries produced during an operation.
type Report struct {
	GeneratedAt time.Time `json:"generated_at" yaml:"generated_at"`
	Source      string    `json:"source"       yaml:"source"`
	Target      string    `json:"target"       yaml:"target"`
	Entries     []Entry   `json:"entries"      yaml:"entries"`
	Summary     Summary   `json:"summary"      yaml:"summary"`
}

// Summary aggregates counts of each operation type.
type Summary struct {
	Added   int `json:"added"   yaml:"added"`
	Removed int `json:"removed" yaml:"removed"`
	Changed int `json:"changed" yaml:"changed"`
	Same    int `json:"same"    yaml:"same"`
	Synced  int `json:"synced"  yaml:"synced"`
	Skipped int `json:"skipped" yaml:"skipped"`
}

// New creates a new Report for the given source and target paths.
func New(source, target string) *Report {
	return &Report{
		GeneratedAt: time.Now().UTC(),
		Source:      source,
		Target:      target,
		Entries:     []Entry{},
	}
}

// Add appends an entry to the report and updates the summary counters.
func (r *Report) Add(e Entry) {
	r.Entries = append(r.Entries, e)
	switch e.Operation {
	case OpAdded:
		r.Summary.Added++
	case OpRemoved:
		r.Summary.Removed++
	case OpChanged:
		r.Summary.Changed++
	case OpSame:
		r.Summary.Same++
	case OpSynced:
		r.Summary.Synced++
	case OpSkipped:
		r.Summary.Skipped++
	}
}

// Write serialises the report to w in the requested format.
func (r *Report) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	case FormatYAML:
		return yaml.NewEncoder(w).Encode(r)
	case FormatText, "":
		return r.writeText(w)
	default:
		return fmt.Errorf("reporter: unsupported format %q", format)
	}
}

// writeText renders a human-readable table to w.
func (r *Report) writeText(w io.Writer) error {
	fmt.Fprintf(w, "Report generated: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Source : %s\n", r.Source)
	fmt.Fprintf(w, "Target : %s\n", r.Target)
	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintf(w, "%-30s %-10s\n", "KEY", "OPERATION")
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, e := range r.Entries {
		maskedNote := ""
		if e.Masked {
			maskedNote = " [masked]"
		}
		fmt.Fprintf(w, "%-30s %-10s%s\n", e.Key, e.Operation, maskedNote)
	}
	fmt.Fprintln(w, strings.Repeat("-", 60))
	s := r.Summary
	fmt.Fprintf(w, "Summary — added:%d removed:%d changed:%d same:%d synced:%d skipped:%d\n",
		s.Added, s.Removed, s.Changed, s.Same, s.Synced, s.Skipped)
	return nil
}
