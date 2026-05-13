package grouper

import (
	"strings"
	"testing"
)

func TestFormat_SortedKeys(t *testing.T) {
	env := map[string]string{
		"APP_Z": "last",
		"APP_A": "first",
		"APP_M": "middle",
	}
	groups := Group(env, Options{Prefixes: []string{"APP"}})
	out := Format(groups, nil)
	idxA := strings.Index(out, "APP_A")
	idxM := strings.Index(out, "APP_M")
	idxZ := strings.Index(out, "APP_Z")
	if !(idxA < idxM && idxM < idxZ) {
		t.Error("keys should appear in sorted order within a group")
	}
}

func TestFormat_KeyCount(t *testing.T) {
	env := map[string]string{"DB_A": "1", "DB_B": "2", "DB_C": "3"}
	groups := Group(env, Options{Prefixes: []string{"DB"}})
	out := Format(groups, nil)
	if !strings.Contains(out, "(3 keys)") {
		t.Errorf("expected key count in banner, got:\n%s", out)
	}
}

func TestFormat_EmptyGroups(t *testing.T) {
	out := Format([]Group{}, nil)
	if out != "" {
		t.Errorf("expected empty string for empty groups, got %q", out)
	}
}

func TestSummary_Empty(t *testing.T) {
	s := Summary([]Group{})
	if s != "" {
		t.Errorf("expected empty summary, got %q", s)
	}
}
