package scorer

import (
	"fmt"
	"sort"
	"strings"
)

// Format returns a human-readable report for a scoring Result.
func Format(r Result) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Score: %d/100  Grade: %s\n", r.Score, r.Grade)

	if len(r.Issues) == 0 {
		sb.WriteString("No issues found.\n")
		return sb.String()
	}

	// Sort issues: high → medium → low, then by key.
	sorted := make([]Issue, len(r.Issues))
	copy(sorted, r.Issues)
	sort.Slice(sorted, func(i, j int) bool {
		si := severityRank(sorted[i].Severity)
		sj := severityRank(sorted[j].Severity)
		if si != sj {
			return si < sj
		}
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(&sb, "Issues (%d):\n", len(sorted))
	for _, iss := range sorted {
		fmt.Fprintf(&sb, "  [%s] %s: %s\n", strings.ToUpper(iss.Severity), iss.Key, iss.Message)
	}
	return sb.String()
}

// Summary returns a one-line summary suitable for CLI output.
func Summary(r Result) string {
	return fmt.Sprintf("grade=%s score=%d issues=%d", r.Grade, r.Score, len(r.Issues))
}

func severityRank(s string) int {
	switch s {
	case "high":
		return 0
	case "medium":
		return 1
	default:
		return 2
	}
}
