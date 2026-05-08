// Package trimmer cleans env maps before further processing.
//
// It can:
//   - Strip leading/trailing whitespace from keys and values.
//   - Drop entries whose value is empty or whitespace-only.
//   - Remove a configurable list of key prefixes (e.g. "APP_", "SERVICE_").
//
// Example:
//
//	opts := trimmer.DefaultOptions()
//	opts.RemoveEmpty = true
//	opts.StripPrefixes = []string{"APP_"}
//	cleaned, err := trimmer.Trim(raw, opts)
package trimmer
