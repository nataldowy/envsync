package pipeline

import (
	"fmt"
	"strings"
)

// StepUpperCaseKeys returns a StepFunc that converts all keys to UPPER_CASE.
func StepUpperCaseKeys() StepFunc {
	return func(env Env) error {
		for k, v := range env {
			upper := strings.ToUpper(k)
			if upper != k {
				delete(env, k)
				env[upper] = v
			}
		}
		return nil
	}
}

// StepStripPrefix returns a StepFunc that removes a key prefix when present.
func StepStripPrefix(prefix string) StepFunc {
	return func(env Env) error {
		for k, v := range env {
			if strings.HasPrefix(k, prefix) {
				delete(env, k)
				env[strings.TrimPrefix(k, prefix)] = v
			}
		}
		return nil
	}
}

// StepRequireKeys returns a StepFunc that errors if any required key is absent.
func StepRequireKeys(keys ...string) StepFunc {
	return func(env Env) error {
		for _, k := range keys {
			if _, ok := env[k]; !ok {
				return fmt.Errorf("required key %q is missing", k)
			}
		}
		return nil
	}
}

// StepSetDefaults returns a StepFunc that sets default values for absent keys.
func StepSetDefaults(defaults Env) StepFunc {
	return func(env Env) error {
		for k, v := range defaults {
			if _, exists := env[k]; !exists {
				env[k] = v
			}
		}
		return nil
	}
}

// StepMaskValues returns a StepFunc that replaces values whose keys match
// sensitive patterns with the provided mask string.
func StepMaskValues(isSensitive func(key string) bool, mask string) StepFunc {
	return func(env Env) error {
		for k, v := range env {
			if v != "" && isSensitive(k) {
				env[k] = mask
			}
		}
		return nil
	}
}
