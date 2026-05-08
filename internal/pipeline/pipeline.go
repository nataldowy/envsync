// Package pipeline provides a composable processing pipeline for .env maps.
// Steps are applied in order; any step returning an error halts execution.
package pipeline

import "fmt"

// Env is a map of key→value pairs representing a parsed .env file.
type Env map[string]string

// StepFunc is a single processing step that transforms an Env in-place.
type StepFunc func(env Env) error

// Pipeline holds an ordered list of steps.
type Pipeline struct {
	steps []StepFunc
}

// New creates an empty Pipeline.
func New() *Pipeline {
	return &Pipeline{}
}

// Add appends one or more steps to the pipeline.
func (p *Pipeline) Add(steps ...StepFunc) *Pipeline {
	p.steps = append(p.steps, steps...)
	return p
}

// Run executes all steps sequentially against a copy of env.
// The original map is never mutated.
func (p *Pipeline) Run(env Env) (Env, error) {
	out := make(Env, len(env))
	for k, v := range env {
		out[k] = v
	}
	for i, step := range p.steps {
		if err := step(out); err != nil {
			return nil, fmt.Errorf("pipeline step %d: %w", i, err)
		}
	}
	return out, nil
}
