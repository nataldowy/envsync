// Package profiler provides static analysis of parsed env maps.
//
// It produces a [Profile] struct that summarises:
//   - total number of keys
//   - keys whose names match common sensitive patterns (password, token, …)
//   - keys with empty or blank values
//   - duplicate key names detected in the raw key slice
//   - the longest key name (useful for formatting)
//
// Example usage:
//
//	env, _ := parser.Parse("production.env")
//	profile := profiler.Analyze(env, nil)
//	fmt.Printf("sensitive keys: %v\n", profile.SensitiveKeys)
package profiler
