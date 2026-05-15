package renamer

import (
	"fmt"
	"sort"
	"strings"
)

// FormatResult returns a human-readable summary of a rename Result.
func FormatResult(r Result) string {
	var sb strings.Builder

	keys := make([]string, 0, len(r.Renamed))
	for old := range r.Renamed {
		keys = append(keys, old)
	}
	sort.Strings(keys)

	if len(keys) > 0 {
		sb.WriteString("Renamed:\n")
		for _, old := range keys {
			fmt.Fprintf(&sb, "  %s -> %s\n", old, r.Renamed[old])
		}
	}

	skipped := make([]string, len(r.Skipped))
	copy(skipped, r.Skipped)
	sort.Strings(skipped)

	if len(skipped) > 0 {
		sb.WriteString("Skipped (no match):\n")
		for _, k := range skipped {
			fmt.Fprintf(&sb, "  %s\n", k)
		}
	}

	return sb.String()
}

// Summary returns a one-line summary of the rename operation.
func Summary(r Result) string {
	return fmt.Sprintf("renamed=%d skipped=%d total=%d",
		len(r.Renamed), len(r.Skipped), len(r.Env))
}
