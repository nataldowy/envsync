package pipeline

import (
	"fmt"
	"strings"

	"github.com/user/envsync/internal/redactor"
)

// StepUpperCaseKeys returns a step that upper-cases all keys.
func StepUpperCaseKeys() Step {
	return func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[strings.ToUpper(k)] = v
		}
		return out, nil
	}
}

// StepStripPrefix returns a step that removes a key prefix.
func StepStripPrefix(prefix string) Step {
	return func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[strings.TrimPrefix(k, prefix)] = v
		}
		return out, nil
	}
}

// StepRequireKeys returns a step that errors if any required key is absent.
func StepRequireKeys(keys ...string) Step {
	return func(env map[string]string) (map[string]string, error) {
		for _, k := range keys {
			if _, ok := env[k]; !ok {
				return nil, fmt.Errorf("required key missing: %s", k)
			}
		}
		return env, nil
	}
}

// StepSetDefaults returns a step that fills missing keys with default values.
func StepSetDefaults(defaults map[string]string) Step {
	return func(env map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = v
		}
		for k, v := range defaults {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
		return out, nil
	}
}

// StepMaskValues returns a step that replaces sensitive values using the
// redactor package, leaving non-sensitive values unchanged.
func StepMaskValues(opts redactor.Options) Step {
	r := redactor.New(opts)
	return func(env map[string]string) (map[string]string, error) {
		return r.RedactMap(env), nil
	}
}
