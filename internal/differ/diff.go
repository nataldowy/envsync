package differ

import (
	"fmt"
	"sort"
	"strings"
)

// EntryStatus represents the diff status of a key between two env files.
type EntryStatus string

const (
	StatusAdded   EntryStatus = "added"
	StatusRemoved EntryStatus = "removed"
	StatusChanged EntryStatus = "changed"
	StatusSame    EntryStatus = "same"
)

// DiffEntry holds the diff result for a single key.
type DiffEntry struct {
	Key      string
	Status   EntryStatus
	OldValue string
	NewValue string
}

// Diff computes the difference between a base env map and a target env map.
// Keys present in base but not target are "removed".
// Keys present in target but not base are "added".
// Keys present in both with different values are "changed".
// Keys present in both with the same value are "same".
func Diff(base, target map[string]string) []DiffEntry {
	keys := make(map[string]struct{})
	for k := range base {
		keys[k] = struct{}{}
	}
	for k := range target {
		keys[k] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	result := make([]DiffEntry, 0, len(sortedKeys))
	for _, k := range sortedKeys {
		baseVal, inBase := base[k]
		targetVal, inTarget := target[k]

		switch {
		case inBase && !inTarget:
			result = append(result, DiffEntry{Key: k, Status: StatusRemoved, OldValue: baseVal})
		case !inBase && inTarget:
			result = append(result, DiffEntry{Key: k, Status: StatusAdded, NewValue: targetVal})
		case baseVal != targetVal:
			result = append(result, DiffEntry{Key: k, Status: StatusChanged, OldValue: baseVal, NewValue: targetVal})
		default:
			result = append(result, DiffEntry{Key: k, Status: StatusSame, OldValue: baseVal, NewValue: targetVal})
		}
	}
	return result
}

// Format returns a human-readable diff string, masking secret values when maskSecrets is true.
func Format(entries []DiffEntry, maskSecrets bool) string {
	var sb strings.Builder
	for _, e := range entries {
		oldVal := e.OldValue
		newVal := e.NewValue
		if maskSecrets {
			if oldVal != "" {
				oldVal = mask(oldVal)
			}
			if newVal != "" {
				newVal = mask(newVal)
			}
		}
		switch e.Status {
		case StatusAdded:
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, newVal))
		case StatusRemoved:
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, oldVal))
		case StatusChanged:
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", e.Key, oldVal, newVal))
		case StatusSame:
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, newVal))
		}
	}
	return sb.String()
}

// mask replaces all but the first and last character of a value with asterisks.
func mask(val string) string {
	if len(val) <= 2 {
		return strings.Repeat("*", len(val))
	}
	return string(val[0]) + strings.Repeat("*", len(val)-2) + string(val[len(val)-1])
}
