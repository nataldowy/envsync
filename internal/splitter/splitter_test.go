package splitter_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/splitter"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"CACHE_HOST":  "redis",
		"APP_PORT":    "8080",
		"UNMATCHED":   "value",
	}
}

func TestSplit_BasicBuckets(t *testing.T) {
	rules := []splitter.Rule{
		{Prefix: "DB_", Bucket: "database"},
		{Prefix: "CACHE_", Bucket: "cache"},
		{Prefix: "APP_", Bucket: "app"},
	}
	res, err := splitter.Split(baseEnv(), rules, splitter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Buckets["database"]["HOST"] != "localhost" {
		t.Errorf("expected database HOST=localhost, got %q", res.Buckets["database"]["HOST"])
	}
	if res.Buckets["cache"]["HOST"] != "redis" {
		t.Errorf("expected cache HOST=redis, got %q", res.Buckets["cache"]["HOST"])
	}
	if res.Buckets["app"]["PORT"] != "8080" {
		t.Errorf("expected app PORT=8080, got %q", res.Buckets["app"]["PORT"])
	}
}

func TestSplit_Unmatched_Included(t *testing.T) {
	rules := []splitter.Rule{
		{Prefix: "DB_", Bucket: "database"},
	}
	opts := splitter.DefaultOptions()
	opts.IncludeUnmatched = true
	res, err := splitter.Split(baseEnv(), rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Buckets["_other"]["UNMATCHED"]; !ok {
		t.Error("expected UNMATCHED key in _other bucket")
	}
}

func TestSplit_Unmatched_Excluded(t *testing.T) {
	rules := []splitter.Rule{
		{Prefix: "DB_", Bucket: "database"},
	}
	opts := splitter.DefaultOptions()
	opts.IncludeUnmatched = false
	res, err := splitter.Split(baseEnv(), rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Buckets["_other"]; ok {
		t.Error("_other bucket should not exist when IncludeUnmatched=false")
	}
}

func TestSplit_StripPrefixDisabled(t *testing.T) {
	rules := []splitter.Rule{
		{Prefix: "DB_", Bucket: "database"},
	}
	opts := splitter.DefaultOptions()
	opts.StripPrefix = false
	res, err := splitter.Split(baseEnv(), rules, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Buckets["database"]["DB_HOST"]; !ok {
		t.Error("expected full key DB_HOST when StripPrefix=false")
	}
}

func TestSplit_LongestPrefixWins(t *testing.T) {
	env := map[string]string{
		"DB_REPLICA_HOST": "replica",
		"DB_HOST":         "primary",
	}
	rules := []splitter.Rule{
		{Prefix: "DB_", Bucket: "db"},
		{Prefix: "DB_REPLICA_", Bucket: "replica"},
	}
	res, err := splitter.Split(env, rules, splitter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Buckets["replica"]["HOST"] != "replica" {
		t.Errorf("expected replica bucket HOST=replica, got %q", res.Buckets["replica"]["HOST"])
	}
	if res.Buckets["db"]["HOST"] != "primary" {
		t.Errorf("expected db bucket HOST=primary, got %q", res.Buckets["db"]["HOST"])
	}
}

func TestSplit_NoRules_Error(t *testing.T) {
	_, err := splitter.Split(baseEnv(), nil, splitter.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty rules, got nil")
	}
}

func TestSplit_EmptyPrefix_Error(t *testing.T) {
	rules := []splitter.Rule{
		{Prefix: "", Bucket: "all"},
	}
	_, err := splitter.Split(baseEnv(), rules, splitter.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty prefix, got nil")
	}
}
