// Package renamer bulk-renames keys in an env map using one of three
// strategies:
//
//   - prefix_swap  – replace a key prefix with a new one
//     (e.g. DEV_DB_HOST -> PROD_DB_HOST)
//
//   - suffix_swap  – replace a key suffix with a new one
//     (e.g. DB_HOST_OLD -> DB_HOST_NEW)
//
//   - regex        – apply a regular-expression replacement to every key
//     that matches the pattern
//
// Keys that do not match the active strategy are either kept as-is
// (SkipUnmatched=false, the default) or dropped from the output
// (SkipUnmatched=true).
//
// The original env map is never mutated; Rename always returns a fresh copy.
package renamer
