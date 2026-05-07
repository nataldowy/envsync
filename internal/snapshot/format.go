package snapshot

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatDelta writes a human-readable summary of a Delta to w.
// Sensitive values are masked when maskFn is non-nil.
func FormatDelta(w io.Writer, d Delta, maskFn func(k, v string) string) {
	if maskFn == nil {
		maskFn = func(_, v string) string { return v }
	}

	writeSection := func(header string, keys []string, valFn func(k string) string) {
		if len(keys) == 0 {
			return
		}
		fmt.Fprintln(w, header)
		for _, k := range keys {
			fmt.Fprintf(w, "  %s=%s\n", k, valFn(k))
		}
	}

	added := sortedKeys(d.Added)
	removed := sortedKeys(d.Removed)
	changed := sortedChangedKeys(d.Changed)

	writeSection("[+] Added:", added, func(k string) string {
		return maskFn(k, d.Added[k])
	})
	writeSection("[-] Removed:", removed, func(k string) string {
		return maskFn(k, d.Removed[k])
	})

	if len(changed) > 0 {
		fmt.Fprintln(w, "[~] Changed:")
		for _, k := range changed {
			pair := d.Changed[k]
			old := maskFn(k, pair[0])
			newV := maskFn(k, pair[1])
			fmt.Fprintf(w, "  %s: %s -> %s\n", k, old, newV)
		}
	}

	if len(added)+len(removed)+len(changed) == 0 {
		fmt.Fprintln(w, "No changes detected.")
	}
}

// Summary returns a one-line summary string for a Delta.
func Summary(d Delta) string {
	parts := []string{}
	if n := len(d.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(d.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(d.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	if len(parts) == 0 {
		return "no changes"
	}
	return strings.Join(parts, ", ")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedChangedKeys(m map[string][2]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
