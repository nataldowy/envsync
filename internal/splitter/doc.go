// Package splitter partitions a flat env map into named buckets using
// prefix-based rules.
//
// # Overview
//
// Large .env files often contain variables for multiple services or
// deployment targets mixed together under a common prefix convention,
// e.g. DB_HOST, CACHE_HOST, APP_PORT.  The splitter package lets you
// declare rules that map each prefix to a named bucket and then splits
// the map in one pass.
//
// # Usage
//
//	rules := []splitter.Rule{
//		{Prefix: "DB_",    Bucket: "database"},
//		{Prefix: "CACHE_", Bucket: "cache"},
//	}
//	res, err := splitter.Split(env, rules, splitter.DefaultOptions())
//	// res.Buckets["database"] => map of DB_ keys (prefix stripped)
//	// res.Buckets["_other"]   => remaining keys
package splitter
