// Package normalizer transforms a map[string]string that represents a parsed
// .env file into a canonical, predictable form.
//
// Supported transformations (all opt-in via Options):
//
//   - UpperCaseKeys  – converts every key to upper-case
//   - TrimValues     – strips leading and trailing whitespace from values
//   - StripPrefix    – removes a fixed prefix from keys before other transforms
//   - AddPrefix      – prepends a fixed string to keys after other transforms
//
// The original map passed to Normalize is never mutated; a new map is always
// returned.
package normalizer
