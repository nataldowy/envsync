package converter

import (
	"fmt"
	"strings"
)

// Summary returns a one-line description of what Convert produced.
func Summary(env map[string]string, format Format) string {
	return fmt.Sprintf("converted %d keys to %s format", len(env), format)
}

// ListFormats returns all supported Format values as strings.
func ListFormats() []string {
	return []string{
		string(FormatExport),
		string(FormatDocker),
		string(FormatInline),
		string(FormatMakefile),
	}
}

// FormatHelp returns a human-readable description for a given format.
func FormatHelp(f Format) string {
	switch f {
	case FormatExport:
		return "POSIX shell: export KEY=VALUE (one per line)"
	case FormatDocker:
		return "Docker CLI: --env KEY=VALUE (one per line)"
	case FormatInline:
		return "Inline shell: KEY=VALUE KEY2=VALUE2 (single line)"
	case FormatMakefile:
		return "GNU Make: KEY := VALUE (one per line)"
	default:
		return fmt.Sprintf("unknown format %q", f)
	}
}

// HelpText returns a formatted block listing all formats and their descriptions.
func HelpText() string {
	var sb strings.Builder
	sb.WriteString("Supported formats:\n")
	for _, name := range ListFormats() {
		fmt.Fprintf(&sb, "  %-12s %s\n", name, FormatHelp(Format(name)))
	}
	return sb.String()
}
