package grouper

import (
	"fmt"
	"sort"
	"strings"
)

// Format renders a slice of Groups as a human-readable string.
// Each group is headed by a banner line followed by sorted key=value pairs.
func Format(groups []Group, maskFn func(key, value string) string) string {
	if maskFn == nil {
		maskFn = func(_, v string) string { return v }
	}

	var sb strings.Builder
	for i, g := range groups {
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "[%s] (%d keys)\n", g.Prefix, len(g.Entries))

		keys := make([]string, 0, len(g.Entries))
		for k := range g.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, maskFn(k, g.Entries[k]))
		}
	}
	return sb.String()
}

// Summary returns a one-line summary of all groups.
func Summary(groups []Group) string {
	parts := make([]string, 0, len(groups))
	for _, g := range groups {
		parts = append(parts, fmt.Sprintf("%s:%d", g.Prefix, len(g.Entries)))
	}
	return strings.Join(parts, "  ")
}
