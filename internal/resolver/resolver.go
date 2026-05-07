// Package resolver provides environment variable resolution with support
// for variable interpolation within .env file values (e.g. FOO=${BAR}_suffix).
package resolver

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls resolver behaviour.
type Options struct {
	// MaxDepth limits recursive resolution to prevent infinite loops.
	MaxDepth int
	// AllowMissing suppresses errors for undefined variables, substituting "".
	AllowMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxDepth:     10,
		AllowMissing: false,
	}
}

// Resolve expands variable references in all values of the provided map.
// It mutates and returns the same map for convenience.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	resolved := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expand(v, env, opts, 0)
		if err != nil {
			return nil, fmt.Errorf("resolver: key %q: %w", k, err)
		}
		resolved[k] = expanded
	}
	return resolved, nil
}

func expand(value string, env map[string]string, opts Options, depth int) (string, error) {
	if depth > opts.MaxDepth {
		return "", fmt.Errorf("max interpolation depth %d exceeded", opts.MaxDepth)
	}

	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return ""
		}
		name := extractName(match)
		val, ok := env[name]
		if !ok {
			if opts.AllowMissing {
				return ""
			}
			expandErr = fmt.Errorf("undefined variable %q", name)
			return ""
		}
		expanded, err := expand(val, env, opts, depth+1)
		if err != nil {
			expandErr = err
			return ""
		}
		return expanded
	})
	if expandErr != nil {
		return "", expandErr
	}
	return strings.TrimSpace(result), nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
