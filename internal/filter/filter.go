// Package filter provides utilities for filtering env key-value pairs
// based on prefix, suffix, pattern matching, or key allowlists/denylists.
package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// Prefix keeps only keys that start with the given prefix.
	Prefix string
	// Suffix keeps only keys that end with the given suffix.
	Suffix string
	// Pattern keeps only keys matching the given regular expression.
	Pattern string
	// AllowList keeps only keys in this set (applied after other filters).
	AllowList []string
	// DenyList removes keys in this set (applied after other filters).
	DenyList []string
}

// DefaultOptions returns Options with no filters applied.
func DefaultOptions() Options {
	return Options{}
}

// Filter returns a new map containing only the entries that pass all
// active filter criteria defined in opts.
func Filter(env map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	deny := make(map[string]bool, len(opts.DenyList))
	for _, k := range opts.DenyList {
		deny[k] = true
	}

	allow := make(map[string]bool, len(opts.AllowList))
	for _, k := range opts.AllowList {
		allow[k] = true
	}

	out := make(map[string]string)
	for k, v := range env {
		if opts.Prefix != "" && !strings.HasPrefix(k, opts.Prefix) {
			continue
		}
		if opts.Suffix != "" && !strings.HasSuffix(k, opts.Suffix) {
			continue
		}
		if re != nil && !re.MatchString(k) {
			continue
		}
		if len(allow) > 0 && !allow[k] {
			continue
		}
		if deny[k] {
			continue
		}
		out[k] = v
	}
	return out, nil
}
