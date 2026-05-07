// Package linter provides style and correctness checks for .env files.
//
// It analyses a parsed key-value map together with a line-number index
// (as produced by internal/validator.buildLineIndex) and returns a slice
// of [Issue] values describing any problems found.
//
// Rules are controlled via [Options]; sensible defaults are available
// through [DefaultOptions].
//
// Example:
//
//	env, _ := parser.Parse("staging.env")
//	idx   := validator.BuildLineIndex("staging.env")
//	issues := linter.Lint(env, idx, linter.DefaultOptions())
//	for _, iss := range issues {
//		fmt.Println(iss)
//	}
package linter
