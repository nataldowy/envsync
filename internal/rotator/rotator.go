// Package rotator provides utilities for rotating secret values in .env maps,
// generating new values and tracking what was changed.
package rotator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// sensitivePattern matches keys that typically hold secrets.
var sensitivePattern = regexp.MustCompile(
	`(?i)(secret|password|passwd|token|key|apikey|api_key|auth|credential|private)`,
)

// Options controls rotation behaviour.
type Options struct {
	// ByteLength is the number of random bytes used to generate a new value.
	ByteLength int
	// Prefix is prepended to every generated value.
	Prefix string
	// OnlyKeys, when non-empty, limits rotation to the specified keys.
	OnlyKeys []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{ByteLength: 32}
}

// Result holds the outcome of a rotation run.
type Result struct {
	// Rotated maps key → new value for every key that was changed.
	Rotated map[string]string
	// Skipped lists keys that were evaluated but not rotated.
	Skipped []string
}

// Rotate iterates over env, replaces sensitive values with freshly generated
// random hex strings, and returns the updated map together with a Result.
// The original map is never mutated.
func Rotate(env map[string]string, opts Options) (map[string]string, Result, error) {
	if opts.ByteLength <= 0 {
		opts.ByteLength = 32
	}

	allowSet := make(map[string]struct{}, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		allowSet[k] = struct{}{}
	}

	out := make(map[string]string, len(env))
	res := Result{Rotated: make(map[string]string)}

	for k, v := range env {
		out[k] = v
	}

	for k := range out {
		if len(allowSet) > 0 {
			if _, ok := allowSet[k]; !ok {
				res.Skipped = append(res.Skipped, k)
				continue
			}
		} else if !sensitivePattern.MatchString(k) {
			res.Skipped = append(res.Skipped, k)
			continue
		}

		newVal, err := generateValue(opts.ByteLength, opts.Prefix)
		if err != nil {
			return nil, Result{}, fmt.Errorf("rotator: generate value for %q: %w", k, err)
		}
		out[k] = newVal
		res.Rotated[k] = newVal
	}

	return out, res, nil
}

func generateValue(byteLen int, prefix string) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.TrimRight(prefix+hex.EncodeToString(b), " "), nil
}
