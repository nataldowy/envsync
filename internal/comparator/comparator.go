// Package comparator provides multi-environment .env comparison,
// allowing you to compare a base environment against multiple targets
// and produce a unified summary of differences.
package comparator

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the comparison outcome for one target environment.
type Result struct {
	Name    string
	Missing []string // keys in base but not in target
	Extra   []string // keys in target but not in base
	Changed []string // keys present in both but with different values
}

// Summary returns a human-readable single-line summary for the result.
func (r Result) Summary() string {
	return fmt.Sprintf("%s: missing=%d extra=%d changed=%d",
		r.Name, len(r.Missing), len(r.Extra), len(r.Changed))
}

// Compare compares base against each named target map.
// base and each entry in targets map key -> value.
func Compare(base map[string]string, targets map[string]map[string]string) []Result {
	names := make([]string, 0, len(targets))
	for n := range targets {
		names = append(names, n)
	}
	sort.Strings(names)

	results := make([]Result, 0, len(targets))
	for _, name := range names {
		target := targets[name]
		res := Result{Name: name}

		for k, bv := range base {
			tv, ok := target[k]
			if !ok {
				res.Missing = append(res.Missing, k)
			} else if tv != bv {
				res.Changed = append(res.Changed, k)
			}
		}
		for k := range target {
			if _, ok := base[k]; !ok {
				res.Extra = append(res.Extra, k)
			}
		}

		sort.Strings(res.Missing)
		sort.Strings(res.Extra)
		sort.Strings(res.Changed)
		results = append(results, res)
	}
	return results
}

// Format renders a slice of Results into a human-readable report string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no environments to compare\n"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("[%s]\n", r.Name))
		writeSection(&sb, "missing", r.Missing)
		writeSection(&sb, "extra", r.Extra)
		writeSection(&sb, "changed", r.Changed)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func writeSection(sb *strings.Builder, label string, keys []string) {
	if len(keys) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("  %s:\n", label))
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("    - %s\n", k))
	}
}
