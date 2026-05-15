// Package sorter provides key-ordering strategies for .env maps.
//
// Supported orderings:
//
//   - alpha      – ascending alphabetical (default)
//   - alpha_desc – descending alphabetical
//   - grouped    – plain keys first (or last) with sensitive keys separated
//   - length     – shortest key name first
//
// Example:
//
//	keys := sorter.Sort(env, sorter.Options{Order: sorter.OrderGrouped})
//	for _, k := range keys {
//		fmt.Printf("%s=%s\n", k, env[k])
//	}
package sorter
