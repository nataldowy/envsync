// Package exporter provides functionality for exporting parsed environment
// variable maps to various output formats: .env (dotenv), JSON, and YAML.
//
// Usage:
//
//	// Export to JSON with sensitive keys masked
//	err := exporter.Export(envMap, "output.json", exporter.FormatJSON, []string{"API_KEY", "DB_PASSWORD"})
//
// Supported formats:
//   - FormatDotEnv  — standard KEY=VALUE format
//   - FormatJSON    — pretty-printed JSON object
//   - FormatYAML    — YAML mapping
//
// Values for keys in the maskKeys list are replaced with "***" before writing.
package exporter
