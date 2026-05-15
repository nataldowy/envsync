package renamer_test

import (
	"testing"

	"github.com/user/envsync/internal/renamer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DEV_DB_HOST":  "localhost",
		"DEV_DB_PORT":  "5432",
		"APP_SECRET":   "s3cr3t",
		"FEATURE_FLAG": "true",
	}
}

func TestRename_PrefixSwap(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.OldPrefix = "DEV_"
	opts.NewPrefix = "PROD_"

	r, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["PROD_DB_HOST"] != "localhost" {
		t.Errorf("expected PROD_DB_HOST=localhost, got %q", r.Env["PROD_DB_HOST"])
	}
	if _, ok := r.Env["DEV_DB_HOST"]; ok {
		t.Error("old key DEV_DB_HOST should not exist in output")
	}
	if len(r.Renamed) != 2 {
		t.Errorf("expected 2 renamed, got %d", len(r.Renamed))
	}
}

func TestRename_SuffixSwap(t *testing.T) {
	env := map[string]string{
		"DB_HOST_OLD": "db.old",
		"DB_PORT_OLD": "5432",
		"KEEP":        "yes",
	}
	opts := renamer.Options{
		Strategy:  renamer.StrategySuffixSwap,
		OldSuffix: "_OLD",
		NewSuffix: "_NEW",
	}
	r, err := renamer.Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["DB_HOST_NEW"] != "db.old" {
		t.Errorf("expected DB_HOST_NEW=db.old, got %q", r.Env["DB_HOST_NEW"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "KEEP" {
		t.Errorf("expected KEEP in skipped, got %v", r.Skipped)
	}
}

func TestRename_Regex(t *testing.T) {
	env := map[string]string{
		"OLD_HOST": "h",
		"OLD_PORT": "p",
		"UNRELATED": "u",
	}
	opts := renamer.Options{
		Strategy:    renamer.StrategyRegex,
		Pattern:     `^OLD_`,
		Replacement: "NEW_",
	}
	r, err := renamer.Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Env["NEW_HOST"] != "h" {
		t.Errorf("expected NEW_HOST=h")
	}
	if r.Env["UNRELATED"] != "u" {
		t.Errorf("expected UNRELATED kept")
	}
}

func TestRename_InvalidRegex(t *testing.T) {
	opts := renamer.Options{
		Strategy: renamer.StrategyRegex,
		Pattern:  "[",
	}
	_, err := renamer.Rename(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestRename_SkipUnmatched(t *testing.T) {
	opts := renamer.Options{
		Strategy:      renamer.StrategyPrefixSwap,
		OldPrefix:     "DEV_",
		NewPrefix:     "PROD_",
		SkipUnmatched: true,
	}
	r, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Env["APP_SECRET"]; ok {
		t.Error("APP_SECRET should be dropped when SkipUnmatched=true")
	}
	if len(r.Env) != 2 {
		t.Errorf("expected 2 keys in output, got %d", len(r.Env))
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	opts := renamer.DefaultOptions()
	opts.OldPrefix = "DEV_"
	opts.NewPrefix = "PROD_"
	_, _ = renamer.Rename(env, opts)
	if _, ok := env["DEV_DB_HOST"]; !ok {
		t.Error("original map was mutated")
	}
}
