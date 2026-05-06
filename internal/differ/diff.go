package differ

import (
	"fmt"
	"sort"
	"strings"
)

// Status represents the type of difference for a key.
type Status int

const (
	Added   Status = iota // present in source, missing in target
	Removed               // missing in source, present in target
	Changed               // present in both, values differ
	Same                  // present in both, values equal
)

// Entry represents a single diff entry.
type Entry struct {
	Key       string
	SourceVal string
	TargetVal string
	Status    Status
}

// Diff computes the difference between source and target env maps.
func Diff(source, target map[string]string) []Entry {
	keys := unionKeys(source, target)
	sort.Strings(keys)

	var entries []Entry
	for _, k := range keys {
		sv, inSrc := source[k]
		tv, inTgt := target[k]

		switch {
		case inSrc && !inTgt:
			entries = append(entries, Entry{Key: k, SourceVal: sv, Status: Added})
		case !inSrc && inTgt:
			entries = append(entries, Entry{Key: k, TargetVal: tv, Status: Removed})
		case sv != tv:
			entries = append(entries, Entry{Key: k, SourceVal: sv, TargetVal: tv, Status: Changed})
		default:
			entries = append(entries, Entry{Key: k, SourceVal: sv, TargetVal: tv, Status: Same})
		}
	}
	return entries
}

// Format renders a human-readable diff string with optional secret masking.
func Format(entries []Entry, maskSecrets bool) string {
	var sb strings.Builder
	for _, e := range entries {
		sv := e.SourceVal
		tv := e.TargetVal
		if maskSecrets {
			sv = mask(sv)
			tv = mask(tv)
		}
		switch e.Status {
		case Added:
			sb.WriteString(fmt.Sprintf("+ %s=%s\n", e.Key, sv))
		case Removed:
			sb.WriteString(fmt.Sprintf("- %s=%s\n", e.Key, tv))
		case Changed:
			sb.WriteString(fmt.Sprintf("~ %s: %s -> %s\n", e.Key, tv, sv))
		case Same:
			sb.WriteString(fmt.Sprintf("  %s=%s\n", e.Key, sv))
		}
	}
	return sb.String()
}

func mask(val string) string {
	if len(val) == 0 {
		return val
	}
	return strings.Repeat("*", len(val))
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
