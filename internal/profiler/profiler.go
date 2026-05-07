// Package profiler analyses an env file and produces a summary report
// describing key statistics: total keys, sensitive keys, empty values, etc.
package profiler

import "strings"

// Profile holds statistics about a parsed env map.
type Profile struct {
	TotalKeys     int
	SensitiveKeys []string
	EmptyValues   []string
	DuplicateKeys []string // detected when built from raw lines
	LongestKey    string
}

// sensitivePatterns mirrors the default patterns used by the masker package.
var sensitivePatterns = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"auth", "credential", "private", "key",
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Analyze builds a Profile from an env map and an optional ordered key slice.
// Pass nil for keys to derive order from the map (non-deterministic).
func Analyze(env map[string]string, keys []string) Profile {
	if keys == nil {
		keys = make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
	}

	seen := make(map[string]int)
	var sensitive, empty, duplicates []string
	longest := ""

	for _, k := range keys {
		seen[k]++
		if seen[k] == 2 {
			duplicates = append(duplicates, k)
		}
		if isSensitive(k) {
			sensitive = append(sensitive, k)
		}
		if v, ok := env[k]; ok && strings.TrimSpace(v) == "" {
			empty = append(empty, k)
		}
		if len(k) > len(longest) {
			longest = k
		}
	}

	return Profile{
		TotalKeys:     len(env),
		SensitiveKeys: sensitive,
		EmptyValues:   empty,
		DuplicateKeys: duplicates,
		LongestKey:    longest,
	}
}
