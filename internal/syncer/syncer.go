package syncer

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/user/envsync/internal/differ"
	"github.com/user/envsync/internal/parser"
)

// SyncMode controls how missing keys are handled during sync.
type SyncMode int

const (
	// ModeAddMissing adds keys present in source but missing in target.
	ModeAddMissing SyncMode = iota
	// ModeOverwrite adds missing keys and updates changed keys.
	ModeOverwrite
)

// Result holds the outcome of a sync operation.
type Result struct {
	Added   []string
	Updated []string
	Skipped []string
}

// HasChanges reports whether the sync result contains any added or updated keys.
func (r *Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Updated) > 0
}

// Sync applies changes from source env map to the target file.
// It reads the target file, applies the diff, and writes the result back.
func Sync(sourceFile, targetFile string, mode SyncMode) (*Result, error) {
	source, err := parser.Parse(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}

	target, err := parser.Parse(targetFile)
	if err != nil {
		return nil, fmt.Errorf("parsing target: %w", err)
	}

	diffs := differ.Diff(source, target)
	result := &Result{}

	for _, d := range diffs {
		switch d.Status {
		case differ.Added:
			// Key in source but not in target — add it.
			target[d.Key] = d.SourceVal
			result.Added = append(result.Added, d.Key)
		case differ.Changed:
			if mode == ModeOverwrite {
				target[d.Key] = d.SourceVal
				result.Updated = append(result.Updated, d.Key)
			} else {
				result.Skipped = append(result.Skipped, d.Key)
			}
		}
	}

	if err := writeEnvFile(targetFile, target); err != nil {
		return nil, fmt.Errorf("writing target: %w", err)
	}

	return result, nil
}

// writeEnvFile serializes an env map to a file in KEY=VALUE format.
func writeEnvFile(path string, env map[string]string) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		if strings.ContainsAny(v, " \t#") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
