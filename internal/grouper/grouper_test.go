package grouper

import (
	"strings"
	"testing"
)

func makeEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"AWS_KEY":     "AKID",
		"AWS_SECRET":  "secret",
		"APP_NAME":    "envsync",
		"LOG_LEVEL":   "info",
		"UNMATCHED":   "yes",
	}
}

func TestGroup_Prefixes(t *testing.T) {
	groups := Group(makeEnv(), Options{Prefixes: []string{"DB", "AWS", "APP"}})
	index := make(map[string]Group)
	for _, g := range groups {
		index[g.Prefix] = g
	}
	if len(index["DB"].Entries) != 2 {
		t.Fatalf("expected 2 DB entries, got %d", len(index["DB"].Entries))
	}
	if len(index["AWS"].Entries) != 2 {
		t.Fatalf("expected 2 AWS entries, got %d", len(index["AWS"].Entries))
	}
	if len(index["APP"].Entries) != 1 {
		t.Fatalf("expected 1 APP entry, got %d", len(index["APP"].Entries))
	}
}

func TestGroup_OtherBucket(t *testing.T) {
	groups := Group(makeEnv(), Options{Prefixes: []string{"DB", "AWS", "APP"}})
	var other *Group
	for i := range groups {
		if groups[i].Prefix == "OTHER" {
			other = &groups[i]
		}
	}
	if other == nil {
		t.Fatal("expected OTHER group")
	}
	if _, ok := other.Entries["LOG_LEVEL"]; !ok {
		t.Error("LOG_LEVEL should be in OTHER")
	}
	if _, ok := other.Entries["UNMATCHED"]; !ok {
		t.Error("UNMATCHED should be in OTHER")
	}
}

func TestGroup_OtherIsLast(t *testing.T) {
	groups := Group(makeEnv(), Options{Prefixes: []string{"DB", "AWS", "APP"}})
	if groups[len(groups)-1].Prefix != "OTHER" {
		t.Error("OTHER should be the last group")
	}
}

func TestGroup_EmptyEnv(t *testing.T) {
	groups := Group(map[string]string{}, Options{Prefixes: []string{"DB"}})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroup_NoPrefixes(t *testing.T) {
	groups := Group(makeEnv(), Options{})
	if len(groups) != 1 || groups[0].Prefix != "OTHER" {
		t.Error("all keys should fall into OTHER when no prefixes given")
	}
}

func TestFormat_ContainsBanner(t *testing.T) {
	groups := Group(makeEnv(), Options{Prefixes: []string{"DB"}})
	out := Format(groups, nil)
	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] banner in output")
	}
}

func TestFormat_MaskFn(t *testing.T) {
	groups := Group(map[string]string{"DB_PASSWORD": "secret"}, Options{Prefixes: []string{"DB"}})
	out := Format(groups, func(k, _ string) string {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") {
			return "***"
		}
		return "visible"
	})
	if !strings.Contains(out, "***") {
		t.Error("expected masked value in output")
	}
}

func TestSummary(t *testing.T) {
	groups := Group(makeEnv(), Options{Prefixes: []string{"DB", "AWS"}})
	s := Summary(groups)
	if !strings.Contains(s, "DB:2") {
		t.Errorf("summary missing DB:2, got: %s", s)
	}
	if !strings.Contains(s, "AWS:2") {
		t.Errorf("summary missing AWS:2, got: %s", s)
	}
}
