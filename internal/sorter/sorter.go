// Package sorter provides utilities for sorting .env key-value maps
// into ordered slices, supporting alphabetical, grouped, and custom orderings.
package sorter

import (
	"regexp"
	"sort"
	"strings"
)

// Order defines the sort strategy.
type Order string

const (
	OrderAlpha      Order = "alpha"       // A-Z by key
	OrderAlphaDesc  Order = "alpha_desc"  // Z-A by key
	OrderGrouped    Order = "grouped"     // sensitive keys last
	OrderLength     Order = "length"      // shortest key first
)

// Options configures sorting behaviour.
type Options struct {
	Order          Order
	SensitiveFirst bool // when true, sensitive keys sort before others in grouped mode
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Order: OrderAlpha}
}

var sensitivePattern = regexp.MustCompile(
	`(?i)(secret|password|passwd|token|api_?key|private|credential|auth)`,
)

// Sort returns the keys of env sorted according to opts.
func Sort(env map[string]string, opts Options) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch opts.Order {
	case OrderAlphaDesc:
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) > strings.ToLower(keys[j])
		})
	case OrderGrouped:
		sort.Slice(keys, func(i, j int) bool {
			si := sensitivePattern.MatchString(keys[i])
			sj := sensitivePattern.MatchString(keys[j])
			if si != sj {
				if opts.SensitiveFirst {
					return si
				}
				return !si
			}
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	case OrderLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) != len(keys[j]) {
				return len(keys[i]) < len(keys[j])
			}
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	default: // OrderAlpha
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	}

	return keys
}
