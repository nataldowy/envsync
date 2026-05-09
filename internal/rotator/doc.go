// Package rotator replaces sensitive secret values inside an env map with
// freshly generated cryptographically random strings.
//
// Usage:
//
//	env := map[string]string{
//		"DB_PASSWORD": "old-secret",
//		"APP_NAME":    "myapp",
//	}
//
//	updated, result, err := rotator.Rotate(env, rotator.DefaultOptions())
//	// updated["DB_PASSWORD"] is a new random hex string
//	// updated["APP_NAME"] is unchanged
//	// result.Rotated contains every key that was replaced
//	// result.Skipped contains every key that was left alone
package rotator
