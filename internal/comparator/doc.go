// Package comparator compares a base .env map against one or more target
// environment maps and reports keys that are missing, extra, or changed.
//
// Typical usage:
//
//	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
//	targets := map[string]map[string]string{
//		"staging": {"HOST": "staging.example.com", "PORT": "5432"},
//		"prod":    {"HOST": "prod.example.com"},
//	}
//	results := comparator.Compare(base, targets)
//	fmt.Print(comparator.Format(results))
package comparator
