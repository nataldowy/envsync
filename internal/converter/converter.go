// Package converter transforms env maps between common configuration formats.
package converter

import (
	"fmt"
	"sort"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatExport  Format = "export"  // export KEY=VALUE shell syntax
	FormatDocker  Format = "docker"  // --env KEY=VALUE flags
	FormatInline  Format = "inline"  // KEY=VALUE KEY2=VALUE2 single line
	FormatMakefile Format = "makefile" // KEY := VALUE
)

// Options controls conversion behaviour.
type Options struct {
	// SortKeys ensures deterministic output.
	SortKeys bool
	// QuoteValues wraps values in double quotes.
	QuoteValues bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		SortKeys:    true,
		QuoteValues: false,
	}
}

// Convert renders env into the requested format string.
func Convert(env map[string]string, format Format, opts Options) (string, error) {
	keys := sortedKeys(env, opts.SortKeys)

	switch format {
	case FormatExport:
		return renderExport(env, keys, opts), nil
	case FormatDocker:
		return renderDocker(env, keys, opts), nil
	case FormatInline:
		return renderInline(env, keys, opts), nil
	case FormatMakefile:
		return renderMakefile(env, keys, opts), nil
	default:
		return "", fmt.Errorf("converter: unsupported format %q", format)
	}
}

func renderExport(env map[string]string, keys []string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%s\n", k, quote(env[k], opts.QuoteValues))
	}
	return sb.String()
}

func renderDocker(env map[string]string, keys []string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "--env %s=%s\n", k, quote(env[k], opts.QuoteValues))
	}
	return sb.String()
}

func renderInline(env map[string]string, keys []string, opts Options) string {
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, quote(env[k], opts.QuoteValues)))
	}
	return strings.Join(parts, " ")
}

func renderMakefile(env map[string]string, keys []string, opts Options) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s := %s\n", k, quote(env[k], opts.QuoteValues))
	}
	return sb.String()
}

func quote(v string, q bool) string {
	if q {
		return `"` + v + `"`
	}
	return v
}

func sortedKeys(env map[string]string, doSort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if doSort {
		sort.Strings(keys)
	}
	return keys
}
