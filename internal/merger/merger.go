// Package merger provides functionality to merge two env file maps,
// with configurable conflict resolution strategies.
package merger

import "fmt"

// Strategy defines how conflicts are resolved when a key exists in both base and override.
type Strategy int

const (
	// StrategyKeepBase keeps the value from the base map on conflict.
	StrategyKeepBase Strategy = iota
	// StrategyTakeOverride replaces the base value with the override value on conflict.
	StrategyTakeOverride
	// StrategyError returns an error on any conflicting key.
	StrategyError
)

// Options configures the merge behaviour.
type Options struct {
	Strategy      Strategy
	SkipEmpty     bool // skip keys with empty values in the override
}

// DefaultOptions returns sensible merge defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:  StrategyTakeOverride,
		SkipEmpty: false,
	}
}

// Result holds the merged map along with metadata about the operation.
type Result struct {
	Merged    map[string]string
	Added     []string // keys that existed only in override
	Overridden []string // keys that existed in both and were overwritten
	Skipped   []string // keys skipped due to empty value or KeepBase strategy
}

// Merge combines base and override maps according to opts.
// The original maps are never mutated.
func Merge(base, override map[string]string, opts Options) (Result, error) {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var result Result

	for k, v := range override {
		if opts.SkipEmpty && v == "" {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if _, exists := merged[k]; exists {
			switch opts.Strategy {
			case StrategyError:
				return Result{}, fmt.Errorf("merger: conflict on key %q", k)
			case StrategyKeepBase:
				result.Skipped = append(result.Skipped, k)
				continue
			case StrategyTakeOverride:
				result.Overridden = append(result.Overridden, k)
			}
		} else {
			result.Added = append(result.Added, k)
		}

		merged[k] = v
	}

	result.Merged = merged
	return result, nil
}
