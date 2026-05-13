// Package grouper partitions a flat env-variable map into named groups
// based on configurable key prefixes (e.g. DB_, AWS_, APP_).
//
// Usage:
//
//	groups := grouper.Group(env, grouper.Options{
//		Prefixes: []string{"DB", "AWS", "APP"},
//	})
//	fmt.Print(grouper.Format(groups, nil))
//
// Keys that match no prefix are collected in the synthetic "OTHER" group,
// which is always placed last in the output.
package grouper
