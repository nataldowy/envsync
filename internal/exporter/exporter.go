package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents the output format for exporting env vars.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
)

// Export writes the given env map to the specified file path in the given format.
// If maskKeys is non-empty, values for those keys are replaced with "***".
func Export(env map[string]string, path string, format Format, maskKeys []string) error {
	masked := applyMask(env, maskKeys)

	var data []byte
	var err error

	switch format {
	case FormatJSON:
		data, err = marshalJSON(masked)
	case FormatYAML:
		data, err = marshalYAML(masked)
	case FormatDotEnv:
		data = []byte(marshalDotEnv(masked))
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}

	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func applyMask(env map[string]string, maskKeys []string) map[string]string {
	set := make(map[string]struct{}, len(maskKeys))
	for _, k := range maskKeys {
		set[strings.ToUpper(k)] = struct{}{}
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, ok := set[strings.ToUpper(k)]; ok {
			out[k] = "***"
		} else {
			out[k] = v
		}
	}
	return out
}

func marshalDotEnv(env map[string]string) string {
	keys := sortedKeys(env)
	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func marshalJSON(env map[string]string) ([]byte, error) {
	ordered := make(map[string]string)
	for k, v := range env {
		ordered[k] = v
	}
	return json.MarshalIndent(ordered, "", "  ")
}

func marshalYAML(env map[string]string) ([]byte, error) {
	return yaml.Marshal(env)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
