package snapshot_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/snapshot"
)

func makeDelta() snapshot.Delta {
	return snapshot.Delta{
		Added:   map[string]string{"NEW_KEY": "hello"},
		Removed: map[string]string{"OLD_KEY": "bye"},
		Changed: map[string][2]string{"DB_PASS": {"old", "new"}},
	}
}

func TestFormatDelta_ContainsSections(t *testing.T) {
	d := makeDelta()
	var sb strings.Builder
	snapshot.FormatDelta(&sb, d, nil)
	out := sb.String()

	for _, want := range []string{"[+] Added", "[-] Removed", "[~] Changed", "NEW_KEY", "OLD_KEY", "DB_PASS"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestFormatDelta_MaskFn(t *testing.T) {
	d := snapshot.Delta{
		Changed: map[string][2]string{"SECRET": {"plain", "other"}},
	}
	var sb strings.Builder
	mask := func(k, v string) string {
		if k == "SECRET" {
			return "***"
		}
		return v
	}
	snapshot.FormatDelta(&sb, d, mask)
	out := sb.String()
	if strings.Contains(out, "plain") || strings.Contains(out, "other") {
		t.Error("expected sensitive values to be masked")
	}
	if !strings.Contains(out, "***") {
		t.Error("expected masked placeholder in output")
	}
}

func TestFormatDelta_NoChanges(t *testing.T) {
	d := snapshot.Delta{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
	}
	var sb strings.Builder
	snapshot.FormatDelta(&sb, d, nil)
	if !strings.Contains(sb.String(), "No changes detected") {
		t.Error("expected no-changes message")
	}
}

func TestSummary_AllThree(t *testing.T) {
	d := makeDelta()
	s := snapshot.Summary(d)
	for _, want := range []string{"added", "removed", "changed"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary missing %q: %s", want, s)
		}
	}
}

func TestSummary_NoChanges(t *testing.T) {
	d := snapshot.Delta{
		Added:   map[string]string{},
		Removed: map[string]string{},
		Changed: map[string][2]string{},
	}
	if got := snapshot.Summary(d); got != "no changes" {
		t.Errorf("got %q want %q", got, "no changes")
	}
}
