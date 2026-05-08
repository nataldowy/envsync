package merger

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMerge_AddsMissingKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"B": "2"}

	res, err := Merge(base, override, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["B"] != "2" {
		t.Errorf("expected B=2, got %q", res.Merged["B"])
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("expected Added=[B], got %v", res.Added)
	}
}

func TestMerge_TakeOverride(t *testing.T) {
	base := map[string]string{"A": "old"}
	override := map[string]string{"A": "new"}

	opts := DefaultOptions()
	opts.Strategy = StrategyTakeOverride
	res, err := Merge(base, override, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "new" {
		t.Errorf("expected A=new, got %q", res.Merged["A"])
	}
	if len(res.Overridden) != 1 || res.Overridden[0] != "A" {
		t.Errorf("expected Overridden=[A], got %v", res.Overridden)
	}
}

func TestMerge_KeepBase(t *testing.T) {
	base := map[string]string{"A": "original"}
	override := map[string]string{"A": "ignored"}

	opts := DefaultOptions()
	opts.Strategy = StrategyKeepBase
	res, err := Merge(base, override, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["A"] != "original" {
		t.Errorf("expected A=original, got %q", res.Merged["A"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected Skipped=[A], got %v", res.Skipped)
	}
}

func TestMerge_StrategyError(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"A": "2"}

	opts := DefaultOptions()
	opts.Strategy = StrategyError
	_, err := Merge(base, override, opts)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_SkipEmpty(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"B": ""}

	opts := DefaultOptions()
	opts.SkipEmpty = true
	res, err := Merge(base, override, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Merged["B"]; ok {
		t.Error("expected B to be skipped")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected Skipped=[B], got %v", res.Skipped)
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	base := map[string]string{"A": "1"}
	override := map[string]string{"A": "2", "B": "3"}

	baseCopy := map[string]string{"A": "1"}

	_, _ = Merge(base, override, DefaultOptions())

	if diff := cmp.Diff(baseCopy, base); diff != "" {
		t.Errorf("base was mutated:\n%s", diff)
	}
}
