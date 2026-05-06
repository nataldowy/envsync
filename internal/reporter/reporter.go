package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourorg/envsync/internal/differ"
)

// Format controls the output format of the reporter.
type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
)

// Stats holds a summary of diff results.
type Stats struct {
	Added   int
	Removed int
	Changed int
	Same    int
}

// Reporter writes diff entries to an io.Writer in a given format.
type Reporter struct {
	w      io.Writer
	fmt    Format
	stats  Stats
}

// New creates a new Reporter writing to w with the given format.
func New(w io.Writer, format Format) *Reporter {
	return &Reporter{w: w, fmt: format}
}

// Write renders all diff entries and accumulates summary stats.
func (r *Reporter) Write(entries []differ.Entry) error {
	// Reset stats each call.
	r.stats = Stats{}

	for _, e := range entries {
		switch e.Status {
		case differ.Added:
			r.stats.Added++
		case differ.Removed:
			r.stats.Removed++
		case differ.Changed:
			r.stats.Changed++
		case differ.Same:
			r.stats.Same++
		}
	}

	switch r.fmt {
	case FormatJSON:
		return r.writeJSON(entries)
	default:
		return r.writePlain(entries)
	}
}

// Summary returns the accumulated Stats from the last Write call.
func (r *Reporter) Summary() Stats {
	return r.stats
}

func (r *Reporter) writePlain(entries []differ.Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(r.w, "no differences found")
		return err
	}

	for _, e := range entries {
		var line string
		switch e.Status {
		case differ.Added:
			line = fmt.Sprintf("+ %s=%s", e.Key, e.RightVal)
		case differ.Removed:
			line = fmt.Sprintf("- %s=%s", e.Key, e.LeftVal)
		case differ.Changed:
			line = fmt.Sprintf("~ %s: %s -> %s", e.Key, e.LeftVal, e.RightVal)
		case differ.Same:
			line = fmt.Sprintf("  %s=%s", e.Key, e.LeftVal)
		}
		if _, err := fmt.Fprintln(r.w, line); err != nil {
			return err
		}
	}
	return nil
}

type jsonEntry struct {
	Key      string `json:"key"`
	Status   string `json:"status"`
	LeftVal  string `json:"left_value,omitempty"`
	RightVal string `json:"right_value,omitempty"`
}

func (r *Reporter) writeJSON(entries []differ.Entry) error {
	out := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, jsonEntry{
			Key:      e.Key,
			Status:   string(e.Status),
			LeftVal:  e.LeftVal,
			RightVal: e.RightVal,
		})
	}
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
