// Package classifier categorises .env keys into logical groups
// (database, auth, network, storage, etc.) based on pattern matching.
package classifier

import (
	"regexp"
	"sort"
	"strings"
)

// Category represents a named group of env keys.
type Category string

const (
	CategoryDatabase Category = "database"
	CategoryAuth     Category = "auth"
	CategoryNetwork  Category = "network"
	CategoryStorage  Category = "storage"
	CategoryLogging  Category = "logging"
	CategoryOther    Category = "other"
)

// rule maps a compiled pattern to a category.
type rule struct {
	pattern  *regexp.Regexp
	category Category
}

// Classifier holds the ordered set of classification rules.
type Classifier struct {
	rules []rule
}

// defaultRules returns the built-in classification rules.
func defaultRules() []rule {
	defs := []struct {
		pattern  string
		category Category
	}{
		{`(?i)(DB_|DATABASE_|POSTGRES|MYSQL|MONGO|REDIS|DSN)`, CategoryDatabase},
		{`(?i)(SECRET|TOKEN|PASSWORD|PASSWD|API_KEY|AUTH|JWT|OAUTH)`, CategoryAuth},
		{`(?i)(HOST|PORT|URL|ENDPOINT|ADDR|DOMAIN|PROXY|TLS|SSL)`, CategoryNetwork},
		{`(?i)(BUCKET|S3|GCS|BLOB|STORAGE|DISK|PATH|DIR|FILE)`, CategoryStorage},
		{`(?i)(LOG|DEBUG|VERBOSE|TRACE|SENTRY|DATADOG)`, CategoryLogging},
	}
	out := make([]rule, 0, len(defs))
	for _, d := range defs {
		out = append(out, rule{pattern: regexp.MustCompile(d.pattern), category: d.category})
	}
	return out
}

// New returns a Classifier using the default built-in rules.
func New() *Classifier {
	return &Classifier{rules: defaultRules()}
}

// Classify returns the Category for a single key.
func (c *Classifier) Classify(key string) Category {
	upper := strings.ToUpper(key)
	for _, r := range c.rules {
		if r.pattern.MatchString(upper) {
			return r.category
		}
	}
	return CategoryOther
}

// ClassifyMap groups all keys in env into categories.
// The returned map is keyed by Category; each value is a sorted slice of keys.
func (c *Classifier) ClassifyMap(env map[string]string) map[Category][]string {
	result := make(map[Category][]string)
	for k := range env {
		cat := c.Classify(k)
		result[cat] = append(result[cat], k)
	}
	for cat := range result {
		sort.Strings(result[cat])
	}
	return result
}
