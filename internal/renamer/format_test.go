package renamer_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/renamer"
)

func TestFormatResult_ContainsRenamedSection(t *testing.T) {
	r := renamer.Result{
		Env:     map[string]string{"PROD_HOST": "h"},
		Renamed: map[string]string{"DEV_HOST": "PROD_HOST"},
		Skipped: []string{},
	}
	out := renamer.FormatResult(r)
	if !strings.Contains(out, "DEV_HOST -> PROD_HOST") {
		t.Errorf("expected rename line in output, got:\n%s", out)
	}
}

func TestFormatResult_ContainsSkippedSection(t *testing.T) {
	r := renamer.Result{
		Env:     map[string]string{"KEEP": "v"},
		Renamed: map[string]string{},
		Skipped: []string{"KEEP"},
	}
	out := renamer.FormatResult(r)
	if !strings.Contains(out, "KEEP") {
		t.Errorf("expected KEEP in skipped section, got:\n%s", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	r := renamer.Result{
		Env:     map[string]string{"A": "1", "B": "2", "C": "3"},
		Renamed: map[string]string{"OLD_A": "A", "OLD_B": "B"},
		Skipped: []string{"C"},
	}
	s := renamer.Summary(r)
	if !strings.Contains(s, "renamed=2") {
		t.Errorf("expected renamed=2 in summary: %s", s)
	}
	if !strings.Contains(s, "skipped=1") {
		t.Errorf("expected skipped=1 in summary: %s", s)
	}
	if !strings.Contains(s, "total=3") {
		t.Errorf("expected total=3 in summary: %s", s)
	}
}

func TestFormatResult_EmptyResult(t *testing.T) {
	r := renamer.Result{
		Env:     map[string]string{},
		Renamed: map[string]string{},
		Skipped: []string{},
	}
	out := renamer.FormatResult(r)
	if out != "" {
		t.Errorf("expected empty output for empty result, got: %q", out)
	}
}
