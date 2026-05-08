// Package pipeline provides a lightweight, composable processing pipeline
// for transforming Env maps (map[string]string) through a series of named
// steps.
//
// Usage:
//
//	p := pipeline.New().
//		Add(pipeline.StepUpperCaseKeys()).
//		Add(pipeline.StepRequireKeys("APP_ENV", "DATABASE_URL")).
//		Add(pipeline.StepSetDefaults(pipeline.Env{"LOG_LEVEL": "info"}))
//
//	out, err := p.Run(inputEnv)
//
// Custom steps implement the StepFunc signature:
//
//	type StepFunc func(env Env) error
package pipeline
