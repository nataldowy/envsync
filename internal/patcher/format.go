package patcher

import (
	"fmt"
	"sort"
	"strings"
)

// FormatResult returns a human-readable summary of the patch result.
func FormatResult(r Result, maskFn func(key, value string) string) string {
	if len(r.Changes) == 0 {
		return "No changes applied.\n"
	}

	if maskFn == nil {
		maskFn = func(_, v string) string { return v }
	}

	// stable order
	changes := make([]Change, len(r.Changes))
	copy(changes, r.Changes)
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Key != changes[j].Key {
			return changes[i].Key < changes[j].Key
		}
		return string(changes[i].Op) < string(changes[j].Op)
	})

	var sb strings.Builder
	for _, c := range changes {
		switch c.Op {
		case OpSet:
			sb.WriteString(fmt.Sprintf("  SET    %s = %s\n", c.Key, maskFn(c.Key, c.NewValue)))
		case OpDelete:
			sb.WriteString(fmt.Sprintf("  DELETE %s (was %s)\n", c.Key, maskFn(c.Key, c.OldValue)))
		case OpRename:
			sb.WriteString(fmt.Sprintf("  RENAME %s -> %s\n", c.Key, c.NewValue))
		}
	}
	return sb.String()
}

// Summary returns a one-line count summary.
func Summary(r Result) string {
	var sets, deletes, renames int
	for _, c := range r.Changes {
		switch c.Op {
		case OpSet:
			sets++
		case OpDelete:
			deletes++
		case OpRename:
			renames++
		}
	}
	return fmt.Sprintf("%d set, %d deleted, %d renamed", sets, deletes, renames)
}
