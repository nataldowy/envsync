// Package sanitizer provides utilities for cleaning and normalising env maps
// before further processing — trimming whitespace, removing banned characters
// from keys, and optionally enforcing upper-case key names.
package sanitizer

import (
	"errors"
	"regexp"
	"strings"
)

// Options controls the behaviour of Sanitize.
type Options struct {
	// UpperCaseKeys converts every key to upper-case when true.
	UpperCaseKeys bool
	// StripInvalidKeyChars removes any character from a key that does not match
	// [A-Za-z0-9_]. Keys that become empty after stripping are dropped.
	StripInvalidKeyChars bool
	// TrimValues trims leading and trailing whitespace from values.
	TrimValues bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		UpperCaseKeys:        true,
		StripInvalidKeyChars: true,
		TrimValues:           true,
	}
}

var invalidKeyChar = regexp.MustCompile(`[^A-Za-z0-9_]`)

// Sanitize applies the given options to src and returns a new map.
// It never mutates src. An error is returned only when StripInvalidKeyChars is
// false and a key contains an invalid character.
func Sanitize(src map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(src))

	for k, v := range src {
		key := k

		if opts.StripInvalidKeyChars {
			key = invalidKeyChar.ReplaceAllString(key, "")
			if key == "" {
				// Key became empty — skip it.
				continue
			}
		} else if invalidKeyChar.MatchString(key) {
			return nil, errors.New("sanitizer: key " + strconv.Quote(k) + " contains invalid characters")
		}

		if opts.UpperCaseKeys {
			key = strings.ToUpper(key)
		}

		val := v
		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}

		out[key] = val
	}

	return out, nil
}
