package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single validation problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
}

func (i Issue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("line %d: %s: %s", i.Line, i.Key, i.Message)
	}
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validate checks a parsed env map (key→value) against a set of rules and
// returns any issues found. rawLines is the original file lines used for
// line-number reporting; pass nil to skip line numbers.
func Validate(env map[string]string, rawLines []string, rules Rules) []Issue {
	var issues []Issue

	lineOf := buildLineIndex(rawLines)

	for key, value := range env {
		line := lineOf[key]

		if !validKeyRe.MatchString(key) {
			issues = append(issues, Issue{Line: line, Key: key, Message: "invalid key name (must match [A-Za-z_][A-Za-z0-9_]*)"})
		}

		if rules.NoEmptyValues && strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{Line: line, Key: key, Message: "empty value"})
		}

		if rules.RequiredKeys != nil {
			// handled separately below
		}
	}

	for _, req := range rules.RequiredKeys {
		if _, ok := env[req]; !ok {
			issues = append(issues, Issue{Key: req, Message: "required key is missing"})
		}
	}

	return issues
}

// Rules configures which validations are performed.
type Rules struct {
	NoEmptyValues bool
	RequiredKeys  []string
}

// buildLineIndex scans raw file lines and returns a map of key→line number.
func buildLineIndex(lines []string) map[string]int {
	index := make(map[string]int)
	if lines == nil {
		return index
	}
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			index[key] = i + 1
		}
	}
	return index
}
