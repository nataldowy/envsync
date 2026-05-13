// Package grouper organises a flat map of env variables into named groups
// based on key prefixes (e.g. DB_, AWS_, APP_).
package grouper

import (
	"sort"
	"strings"
)

// Group holds all key/value pairs that share a common prefix.
type Group struct {
	Prefix string
	Entries map[string]string
}

// Options controls how grouping is performed.
type Options struct {
	// Prefixes is the ordered list of prefixes to match.
	// Keys that match no prefix are placed in the "OTHER" group.
	Prefixes []string
	// Separator is the character that terminates a prefix (default "_").
	Separator string
	// CaseSensitive controls whether prefix matching is case-sensitive.
	CaseSensitive bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Separator:     "_",
		CaseSensitive: false,
	}
}

// Group partitions env into named groups according to opts.
// The returned slice is sorted by prefix name; "OTHER" is always last.
func Group(env map[string]string, opts Options) []Group {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	index := make(map[string]*Group)

	for k, v := range env {
		matched := false
		for _, p := range opts.Prefixes {
			candidate := p + opts.Separator
			key := k
			if !opts.CaseSensitive {
				candidate = strings.ToUpper(candidate)
				key = strings.ToUpper(k)
			}
			if strings.HasPrefix(key, candidate) {
				norm := strings.ToUpper(p)
				if _, ok := index[norm]; !ok {
					index[norm] = &Group{Prefix: norm, Entries: make(map[string]string)}
				}
				index[norm].Entries[k] = v
				matched = true
				break
			}
		}
		if !matched {
			if _, ok := index["OTHER"]; !ok {
				index["OTHER"] = &Group{Prefix: "OTHER", Entries: make(map[string]string)}
			}
			index["OTHER"].Entries[k] = v
		}
	}

	result := make([]Group, 0, len(index))
	for _, g := range index {
		result = append(result, *g)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Prefix == "OTHER" {
			return false
		}
		if result[j].Prefix == "OTHER" {
			return true
		}
		return result[i].Prefix < result[j].Prefix
	})
	return result
}
