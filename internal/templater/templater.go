// Package templater renders .env files from a template with placeholder substitution.
// Templates use {{KEY}} syntax; values are sourced from a provided map.
package templater

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Options controls templater behaviour.
type Options struct {
	// FailOnMissing causes Render to return an error if a placeholder has no value.
	FailOnMissing bool
	// FillEmpty replaces missing placeholders with an empty string instead of
	// leaving the original token when FailOnMissing is false.
	FillEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: true,
		FillEmpty:     false,
	}
}

var placeholderRe = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

// Render substitutes {{KEY}} placeholders in tmpl with values from vars.
// Lines beginning with '#' are passed through unchanged.
func Render(tmpl string, vars map[string]string, opts Options) (string, error) {
	var sb strings.Builder
	lines := strings.Split(tmpl, "\n")

	for i, line := range lines {
		resolved, err := renderLine(line, vars, opts)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", i+1, err)
		}
		sb.WriteString(resolved)
		if i < len(lines)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// Placeholders returns all unique placeholder names found in tmpl, sorted.
func Placeholders(tmpl string) []string {
	seen := map[string]struct{}{}
	for _, m := range placeholderRe.FindAllStringSubmatch(tmpl, -1) {
		seen[m[1]] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func renderLine(line string, vars map[string]string, opts Options) (string, error) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "#") {
		return line, nil
	}
	var rerr error
	result := placeholderRe.ReplaceAllStringFunc(line, func(match string) string {
		if rerr != nil {
			return match
		}
		key := placeholderRe.FindStringSubmatch(match)[1]
		val, ok := vars[key]
		if !ok {
			if opts.FailOnMissing {
				rerr = fmt.Errorf("missing value for placeholder %q", key)
				return match
			}
			if opts.FillEmpty {
				return ""
			}
			return match
		}
		return val
	})
	return result, rerr
}
