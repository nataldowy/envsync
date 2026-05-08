// Package trimmer provides utilities for cleaning up env maps:
// removing keys with blank values, trimming whitespace from keys
// and values, and stripping a configurable set of key prefixes.
package trimmer

import "strings"

// Options controls Trim behaviour.
type Options struct {
	// RemoveEmpty drops keys whose value is empty or whitespace-only.
	RemoveEmpty bool
	// TrimSpace strips leading/trailing whitespace from every key and value.
	TrimSpace bool
	// StripPrefixes removes the first matching prefix from each key.
	StripPrefixes []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		RemoveEmpty: false,
		TrimSpace:   true,
		StripPrefixes: nil,
	}
}

// Trim returns a new map derived from src after applying opts.
// The original map is never mutated.
func Trim(src map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(src))

	for k, v := range src {
		key := k
		val := v

		if opts.TrimSpace {
			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)
		}

		if opts.RemoveEmpty && strings.TrimSpace(val) == "" {
			continue
		}

		for _, pfx := range opts.StripPrefixes {
			if strings.HasPrefix(key, pfx) {
				key = strings.TrimPrefix(key, pfx)
				break
			}
		}

		if key == "" {
			continue
		}

		out[key] = val
	}

	return out, nil
}
