package comparator

import (
	"strings"
	"testing"
)

var base = map[string]string{
	"HOST":     "localhost",
	"PORT":     "5432",
	"DB_NAME":  "mydb",
	"API_KEY":  "secret",
}

func TestCompare_Missing(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"HOST": "localhost", "PORT": "5432", "DB_NAME": "mydb"},
	}
	results := Compare(base, targets)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Missing) != 1 || results[0].Missing[0] != "API_KEY" {
		t.Errorf("expected API_KEY missing, got %v", results[0].Missing)
	}
}

func TestCompare_Extra(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"HOST": "localhost", "PORT": "5432", "DB_NAME": "mydb", "API_KEY": "secret", "NEW_VAR": "val"},
	}
	results := Compare(base, targets)
	if len(results[0].Extra) != 1 || results[0].Extra[0] != "NEW_VAR" {
		t.Errorf("expected NEW_VAR extra, got %v", results[0].Extra)
	}
}

func TestCompare_Changed(t *testing.T) {
	targets := map[string]map[string]string{
		"prod": {"HOST": "prod.host", "PORT": "5432", "DB_NAME": "mydb", "API_KEY": "secret"},
	}
	results := Compare(base, targets)
	if len(results[0].Changed) != 1 || results[0].Changed[0] != "HOST" {
		t.Errorf("expected HOST changed, got %v", results[0].Changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	copy := map[string]string{"HOST": "localhost", "PORT": "5432", "DB_NAME": "mydb", "API_KEY": "secret"}
	results := Compare(base, map[string]map[string]string{"same": copy})
	r := results[0]
	if len(r.Missing)+len(r.Extra)+len(r.Changed) != 0 {
		t.Errorf("expected no differences, got %+v", r)
	}
}

func TestCompare_SortedResults(t *testing.T) {
	targets := map[string]map[string]string{
		"z-env": {},
		"a-env": {},
	}
	results := Compare(base, targets)
	if results[0].Name != "a-env" || results[1].Name != "z-env" {
		t.Errorf("expected sorted result names, got %s %s", results[0].Name, results[1].Name)
	}
}

func TestCompare_EmptyTargets(t *testing.T) {
	results := Compare(base, map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestFormat_ContainsSections(t *testing.T) {
	targets := map[string]map[string]string{
		"staging": {"HOST": "other"},
	}
	results := Compare(base, targets)
	out := Format(results)
	for _, want := range []string{"[staging]", "missing", "changed"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output missing %q", want)
		}
	}
}

func TestFormat_Empty(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "no environments") {
		t.Errorf("expected no-environments message, got %q", out)
	}
}

func TestResult_Summary(t *testing.T) {
	r := Result{Name: "prod", Missing: []string{"A"}, Extra: []string{"B", "C"}, Changed: []string{}}
	s := r.Summary()
	if !strings.Contains(s, "prod") || !strings.Contains(s, "missing=1") || !strings.Contains(s, "extra=2") {
		t.Errorf("unexpected summary: %s", s)
	}
}
