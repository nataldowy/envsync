// Package normalizer provides utilities to normalize .env key-value maps
// into a canonical form: upper-case keys, trimmed values, and optional
// prefix injection.
package normalizer

import (
	"strings"
)

// Options controls the behaviour of Normalize.
type Options struct {
	// UpperCaseKeys converts every key to upper-case when true.
	UpperCaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// AddPrefix prepends the given string to every key (applied after
	// UpperCaseKeys so the prefix is preserved as-is).
	AddPrefix string
	// StripPrefix removes the given prefix from keys before any other
	// transformation.
	StripPrefix string
}

// DefaultOptions returns an Options value with the most common settings.
func DefaultOptions() Options {
	return Options{
		UpperCaseKeys: true,
		TrimValues:    true,
	}
}

// Normalize returns a new map that is the normalised form of env according to
// opts. The original map is never mutated.
func Normalize(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		// 1. strip prefix
		if opts.StripPrefix != "" {
			k = strings.TrimPrefix(k, opts.StripPrefix)
		}
		// 2. upper-case
		if opts.UpperCaseKeys {
			k = strings.ToUpper(k)
		}
		// 3. add prefix
		if opts.AddPrefix != "" {
			k = opts.AddPrefix + k
		}
		// 4. trim value
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		// skip keys that became empty after transformations
		if k == "" {
			continue
		}
		out[k] = v
	}
	return out
}
