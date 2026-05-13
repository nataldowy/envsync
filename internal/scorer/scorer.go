// Package scorer evaluates the overall quality and security score of a .env file
// based on various heuristics such as sensitive key exposure, empty values,
// naming conventions, and duplicate entries.
package scorer

import (
	"regexp"
	"strings"
)

// Result holds the outcome of a scoring evaluation.
type Result struct {
	Score    int      // 0–100
	Grade    string   // A, B, C, D, F
	Issues   []Issue
}

// Issue describes a single problem found during scoring.
type Issue struct {
	Key      string
	Severity string // "high", "medium", "low"
	Message  string
}

var sensitivePattern = regexp.MustCompile(
	`(?i)(password|secret|token|key|api_?key|private|credential|auth)`,
)

// Score evaluates a map of env vars and returns a Result.
func Score(env map[string]string) Result {
	var issues []Issue
	penalty := 0

	seen := make(map[string]int)
	for k := range env {
		seen[strings.ToUpper(k)]++
	}

	for k, v := range env {
		upper := strings.ToUpper(k)

		// Duplicate keys (case-insensitive)
		if seen[upper] > 1 {
			issues = append(issues, Issue{Key: k, Severity: "medium", Message: "duplicate key (case-insensitive)"})
			penalty += 5
		}

		// Sensitive key with empty value
		if sensitivePattern.MatchString(k) && v == "" {
			issues = append(issues, Issue{Key: k, Severity: "high", Message: "sensitive key has empty value"})
			penalty += 15
		}

		// Non-uppercase key name
		if k != strings.ToUpper(k) {
			issues = append(issues, Issue{Key: k, Severity: "low", Message: "key is not upper-case"})
			penalty += 2
		}

		// Plaintext-looking secrets (short alphanumeric values on sensitive keys)
		if sensitivePattern.MatchString(k) && len(v) > 0 && len(v) < 8 {
			issues = append(issues, Issue{Key: k, Severity: "high", Message: "sensitive key has suspiciously short value"})
			penalty += 10
		}
	}

	if penalty > 100 {
		penalty = 100
	}
	score := 100 - penalty

	return Result{
		Score:  score,
		Grade:  grade(score),
		Issues: issues,
	}
}

func grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
