package masker

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a key holds a secret.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"ACCESS_KEY",
}

const MaskedValue = "********"

// Masker decides which env values should be hidden.
type Masker struct {
	patterns []string
}

// New returns a Masker using the provided patterns.
// If patterns is nil, DefaultSensitivePatterns is used.
func New(patterns []string) *Masker {
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return &Masker{patterns: upper}
}

// IsSensitive reports whether the given key name matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range m.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Mask returns the masked placeholder if the key is sensitive, otherwise the
// original value.
func (m *Masker) Mask(key, value string) string {
	if m.IsSensitive(key) {
		return MaskedValue
	}
	return value
}

// MaskMap returns a copy of env where sensitive values are replaced.
func (m *Masker) MaskMap(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = m.Mask(k, v)
	}
	return out
}
