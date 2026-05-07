// Package resolver implements variable interpolation for .env files.
//
// It supports two common reference styles:
//
//	${VAR_NAME}   — brace-delimited reference
//	$VAR_NAME     — bare dollar reference
//
// References are resolved recursively up to a configurable depth limit
// (default 10) to prevent infinite loops caused by circular definitions.
//
// Example:
//
//	env := map[string]string{
//		"BASE_URL": "https://example.com",
//		"API_URL":  "${BASE_URL}/api/v1",
//	}
//	resolved, err := resolver.Resolve(env, resolver.DefaultOptions())
//	// resolved["API_URL"] == "https://example.com/api/v1"
package resolver
