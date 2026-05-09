package rotator

import (
	"fmt"
	"sort"
	"strings"
)

// FormatResult renders a human-readable summary of a rotation Result.
//
// Example output:
//
//	Rotated (2):
//	  DB_PASSWORD  →  <new value hidden>
//	  API_KEY      →  <new value hidden>
//
//	Skipped (1):
//	  APP_NAME
func FormatResult(res Result, showValues bool) string {
	var sb strings.Builder

	rotatedKeys := sortedKeys(res.Rotated)

	fmt.Fprintf(&sb, "Rotated (%d):\n", len(rotatedKeys))
	for _, k := range rotatedKeys {
		if showValues {
			fmt.Fprintf(&sb, "  %s  →  %s\n", k, res.Rotated[k])
		} else {
			fmt.Fprintf(&sb, "  %s  →  <new value hidden>\n", k)
		}
	}

	skipped := make([]string, len(res.Skipped))
	copy(skipped, res.Skipped)
	sort.Strings(skipped)

	fmt.Fprintf(&sb, "\nSkipped (%d):\n", len(skipped))
	for _, k := range skipped {
		fmt.Fprintf(&sb, "  %s\n", k)
	}

	return sb.String()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
