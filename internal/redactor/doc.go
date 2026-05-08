// Package redactor provides sensitive-value redaction for .env content.
//
// It identifies sensitive keys by matching against a set of built-in
// regular expressions (covering common patterns such as "password",
// "secret", "token", "key", etc.) and an optional caller-supplied
// extra pattern.
//
// Three entry points are provided:
//
//   - RedactMap  – operates on a map[string]string, returning a copy.
//   - RedactLine – operates on a single KEY=VALUE text line.
//   - RedactLines – operates on a full multi-line .env string.
//
// Usage:
//
//	r := redactor.New(redactor.Options{Placeholder: "***"})
//	safe := r.RedactMap(env)
package redactor
