// Package splitter splits a flat env map into multiple named buckets
// based on key prefix rules, producing one map per target environment
// or service boundary.
package splitter

import (
	"fmt"
	"strings"
)

// Options controls splitting behaviour.
type Options struct {
	// StripPrefix removes the matched prefix from keys in the output bucket.
	StripPrefix bool
	// IncludeUnmatched places keys that match no rule into the "_other" bucket.
	IncludeUnmatched bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		StripPrefix:      true,
		IncludeUnmatched: true,
	}
}

// Rule maps a prefix to a named output bucket.
type Rule struct {
	Prefix string
	Bucket string
}

// Result holds the split output.
type Result struct {
	Buckets map[string]map[string]string
}

// Split partitions env into named buckets according to rules.
// Keys are matched by longest prefix first to avoid ambiguity.
func Split(env map[string]string, rules []Rule, opts Options) (Result, error) {
	if len(rules) == 0 {
		return Result{}, fmt.Errorf("splitter: at least one rule is required")
	}

	// Validate rules.
	for _, r := range rules {
		if strings.TrimSpace(r.Prefix) == "" {
			return Result{}, fmt.Errorf("splitter: rule bucket %q has an empty prefix", r.Bucket)
		}
		if strings.TrimSpace(r.Bucket) == "" {
			return Result{}, fmt.Errorf("splitter: rule with prefix %q has an empty bucket name", r.Prefix)
		}
	}

	// Sort rules by descending prefix length for longest-match semantics.
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if len(sorted[j].Prefix) > len(sorted[i].Prefix) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	buckets := make(map[string]map[string]string)

	for key, val := range env {
		matched := false
		for _, r := range sorted {
			if strings.HasPrefix(key, r.Prefix) {
				if buckets[r.Bucket] == nil {
					buckets[r.Bucket] = make(map[string]string)
				}
				outKey := key
				if opts.StripPrefix {
					outKey = strings.TrimPrefix(key, r.Prefix)
					if outKey == "" {
						outKey = key // never produce an empty key
					}
				}
				buckets[r.Bucket][outKey] = val
				matched = true
				break
			}
		}
		if !matched && opts.IncludeUnmatched {
			if buckets["_other"] == nil {
				buckets["_other"] = make(map[string]string)
			}
			buckets["_other"][key] = val
		}
	}

	return Result{Buckets: buckets}, nil
}
