// Package linter checks .env files for common style and correctness issues.
package linter

import (
	"fmt"
	"strings"
	"unicode"
)

// Severity indicates how serious a lint issue is.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Issue represents a single lint finding.
type Issue struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] line %d (%s): %s", i.Severity, i.Line, i.Key, i.Message)
}

// Options controls which lint rules are active.
type Options struct {
	AllowEmptyValues   bool
	RequireUpperCase   bool
	MaxKeyLength       int
	ForbidLeadingDigit bool
}

// DefaultOptions returns a sensible default lint configuration.
func DefaultOptions() Options {
	return Options{
		AllowEmptyValues:   true,
		RequireUpperCase:   true,
		MaxKeyLength:       64,
		ForbidLeadingDigit: true,
	}
}

// Lint analyses the provided key-value map (line numbers keyed by variable name)
// and returns any issues found.
func Lint(entries map[string]string, lineIndex map[string]int, opts Options) []Issue {
	var issues []Issue

	for key, value := range entries {
		line := lineIndex[key]

		if opts.RequireUpperCase && key != strings.ToUpper(key) {
			issues = append(issues, Issue{
				Line:     line,
				Key:      key,
				Message:  "key should be UPPER_CASE",
				Severity: SeverityWarning,
			})
		}

		if opts.ForbidLeadingDigit && len(key) > 0 && unicode.IsDigit(rune(key[0])) {
			issues = append(issues, Issue{
				Line:     line,
				Key:      key,
				Message:  "key must not start with a digit",
				Severity: SeverityError,
			})
		}

		if opts.MaxKeyLength > 0 && len(key) > opts.MaxKeyLength {
			issues = append(issues, Issue{
				Line:     line,
				Key:      key,
				Message:  fmt.Sprintf("key length %d exceeds maximum %d", len(key), opts.MaxKeyLength),
				Severity: SeverityWarning,
			})
		}

		if !opts.AllowEmptyValues && strings.TrimSpace(value) == "" {
			issues = append(issues, Issue{
				Line:     line,
				Key:      key,
				Message:  "empty value is not allowed",
				Severity: SeverityError,
			})
		}
	}

	return issues
}
