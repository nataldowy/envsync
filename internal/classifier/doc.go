// Package classifier provides key-level categorisation for .env maps.
//
// Keys are matched against a prioritised list of regular expressions and
// assigned to one of the built-in categories: database, auth, network,
// storage, logging, or other.
//
// Usage:
//
//	c := classifier.New()
//	cat := c.Classify("DB_HOST")          // → classifier.CategoryDatabase
//	groups := c.ClassifyMap(envMap)        // → map[Category][]string
package classifier
