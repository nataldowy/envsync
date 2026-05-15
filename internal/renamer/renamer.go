// Package renamer provides utilities for bulk-renaming keys in an env map
// according to a set of rules (prefix swap, suffix swap, or regex replacement).
package renamer

import (
	"fmt"
	"regexp"
	"strings"
)

// Strategy controls how keys are renamed.
type Strategy string

const (
	StrategyPrefixSwap Strategy = "prefix_swap"
	StrategySuffixSwap Strategy = "suffix_swap"
	StrategyRegex      Strategy = "regex"
)

// Options configures a Rename operation.
type Options struct {
	Strategy  Strategy
	OldPrefix string
	NewPrefix string
	OldSuffix string
	NewSuffix string
	Pattern   string // used with StrategyRegex
	Replacement string
	SkipUnmatched bool // if false, unmatched keys are kept as-is
}

// DefaultOptions returns sensible defaults (prefix swap, no-op).
func DefaultOptions() Options {
	return Options{
		Strategy:      StrategyPrefixSwap,
		SkipUnmatched: false,
	}
}

// Result holds the renamed map and metadata about what changed.
type Result struct {
	Env     map[string]string
	Renamed map[string]string // oldKey -> newKey
	Skipped []string          // keys that did not match and were kept
}

// Rename applies the rename strategy to env and returns a Result.
// The original map is never mutated.
func Rename(env map[string]string, opts Options) (Result, error) {
	renamed := make(map[string]string, len(env))
	out := make(map[string]string, len(env))
	var skipped []string

	var re *regexp.Regexp
	if opts.Strategy == StrategyRegex {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return Result{}, fmt.Errorf("renamer: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	for k, v := range env {
		newKey, matched := applyStrategy(k, opts, re)
		if !matched {
			skipped = append(skipped, k)
			if !opts.SkipUnmatched {
				out[k] = v
			}
			continue
		}
		if existing, conflict := out[newKey]; conflict && existing != v {
			return Result{}, fmt.Errorf("renamer: key collision on %q", newKey)
		}
		out[newKey] = v
		renamed[k] = newKey
	}
	return Result{Env: out, Renamed: renamed, Skipped: skipped}, nil
}

func applyStrategy(key string, opts Options, re *regexp.Regexp) (string, bool) {
	switch opts.Strategy {
	case StrategyPrefixSwap:
		if !strings.HasPrefix(key, opts.OldPrefix) {
			return key, false
		}
		return opts.NewPrefix + strings.TrimPrefix(key, opts.OldPrefix), true
	case StrategySuffixSwap:
		if !strings.HasSuffix(key, opts.OldSuffix) {
			return key, false
		}
		return strings.TrimSuffix(key, opts.OldSuffix) + opts.NewSuffix, true
	case StrategyRegex:
		if re == nil || !re.MatchString(key) {
			return key, false
		}
		return re.ReplaceAllString(key, opts.Replacement), true
	}
	return key, false
}
