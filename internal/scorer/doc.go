// Package scorer computes a quality and security score for a set of environment
// variables. It inspects keys and values for common problems — such as sensitive
// keys with empty or weak values, non-uppercase naming, and duplicates — and
// produces a numeric score (0–100) together with a letter grade and a list of
// actionable issues.
//
// Usage:
//
//	env := map[string]string{"DB_PASSWORD": "", "api_key": "abc"}
//	result := scorer.Score(env)
//	fmt.Println(scorer.Format(result))
package scorer
