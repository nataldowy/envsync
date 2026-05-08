// Package redactor provides line-level redaction of .env file content,
// replacing sensitive values with a configurable placeholder before output.
package redactor

import (
	"fmt"
	"regexp"
	"strings"
)

// DefaultPlaceholder is used when no custom placeholder is set.
const DefaultPlaceholder = "[REDACTED]"

// Options controls redactor behaviour.
type Options struct {
	Placeholder  string
	ExtraPattern *regexp.Regexp
}

// Redactor holds compiled sensitive-key patterns.
type Redactor struct {
	placeholder string
	patterns    []*regexp.Regexp
}

var defaultPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(secret|password|passwd|token|key|apikey|api_key|auth|credential|private)`),
}

// New creates a Redactor with the given options.
func New(opts Options) *Redactor {
	ph := opts.Placeholder
	if ph == "" {
		ph = DefaultPlaceholder
	}
	pats := make([]*regexp.Regexp, len(defaultPatterns))
	copy(pats, defaultPatterns)
	if opts.ExtraPattern != nil {
		pats = append(pats, opts.ExtraPattern)
	}
	return &Redactor{placeholder: ph, patterns: pats}
}

// IsSensitive reports whether the key name matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	for _, p := range r.patterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

// RedactMap returns a copy of env where sensitive values are replaced.
func (r *Redactor) RedactMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactLine redacts the value portion of a KEY=VALUE line if the key is sensitive.
// Comment and blank lines are returned unchanged.
func (r *Redactor) RedactLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return line
	}
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return line
	}
	key := strings.TrimSpace(line[:idx])
	if r.IsSensitive(key) {
		return fmt.Sprintf("%s=%s", key, r.placeholder)
	}
	return line
}

// RedactLines applies RedactLine to every line in src.
func (r *Redactor) RedactLines(src string) string {
	lines := strings.Split(src, "\n")
	for i, l := range lines {
		lines[i] = r.RedactLine(l)
	}
	return strings.Join(lines, "\n")
}
