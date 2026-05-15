// Package converter renders an env map into alternative configuration
// formats such as shell export statements, Docker --env flags, single-line
// inline assignments, and Makefile variable declarations.
//
// Basic usage:
//
//	result, err := converter.Convert(
//		env,
//		converter.FormatExport,
//		converter.DefaultOptions(),
//	)
package converter
