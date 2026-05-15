// Package deduplicator detects and removes duplicate keys from an env map,
// preserving the last-seen value by default or reporting conflicts.
package deduplicator

import (
	"fmt"
	"sort"
	"strings"
)

// Strategy controls which value wins when a duplicate key is found.
type Strategy string

const (
	// StrategyLast keeps the last occurrence of a duplicate key (default).
	StrategyLast Strategy = "last"
	// StrategyFirst keeps the first occurrence of a duplicate key.
	StrategyFirst Strategy = "first"
	// StrategyError returns an error if any duplicate key is found.
	StrategyError Strategy = "error"
)

// Result holds the deduplicated env map and metadata about duplicates found.
type Result struct {
	Env        map[string]string
	Duplicates map[string][]string // key -> all values seen
}

// Options configures deduplication behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{Strategy: StrategyLast}
}

// Deduplicate processes a slice of raw key=value pairs (preserving order) and
// returns a Result according to the chosen strategy. Blank lines and comments
// are silently ignored.
func Deduplicate(lines []string, opts Options) (Result, error) {
	if opts.Strategy == "" {
		opts.Strategy = StrategyLast
	}

	seen := make(map[string][]string) // key -> ordered values
	order := []string{}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if key == "" {
			continue
		}
		if _, exists := seen[key]; !exists {
			order = append(order, key)
		}
		seen[key] = append(seen[key], val)
	}

	// Collect actual duplicates.
	duplicates := make(map[string][]string)
	for k, vals := range seen {
		if len(vals) > 1 {
			duplicates[k] = vals
		}
	}

	if opts.Strategy == StrategyError && len(duplicates) > 0 {
		keys := make([]string, 0, len(duplicates))
		for k := range duplicates {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return Result{}, fmt.Errorf("duplicate keys found: %s", strings.Join(keys, ", "))
	}

	env := make(map[string]string, len(order))
	for _, key := range order {
		vals := seen[key]
		switch opts.Strategy {
		case StrategyFirst:
			env[key] = vals[0]
		default: // StrategyLast
			env[key] = vals[len(vals)-1]
		}
	}

	return Result{Env: env, Duplicates: duplicates}, nil
}
