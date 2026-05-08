package pipeline_test

import (
	"testing"

	"github.com/user/envsync/internal/pipeline"
	"github.com/user/envsync/internal/redactor"
)

func TestStepMaskValues_EmptyValueNotMasked(t *testing.T) {
	step := pipeline.StepMaskValues(redactor.Options{})
	env := map[string]string{"DB_PASSWORD": "", "APP_ENV": "dev"}
	out, err := step(env)
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_PASSWORD"] != redactor.DefaultPlaceholder {
		t.Errorf("expected placeholder for empty sensitive value, got %q", out["DB_PASSWORD"])
	}
	if out["APP_ENV"] != "dev" {
		t.Errorf("plain value mutated: %q", out["APP_ENV"])
	}
}

func TestStepRequireKeys_AllPresent(t *testing.T) {
	step := pipeline.StepRequireKeys("A", "B")
	_, err := step(map[string]string{"A": "1", "B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChainedSteps(t *testing.T) {
	p := pipeline.New(
		pipeline.StepSetDefaults(map[string]string{"PORT": "3000"}),
		pipeline.StepUpperCaseKeys(),
		pipeline.StepMaskValues(redactor.Options{Placeholder: "XXX"}),
	)
	env := map[string]string{"api_key": "secret", "app_env": "staging"}
	out, err := p.Run(env)
	if err != nil {
		t.Fatal(err)
	}
	if out["API_KEY"] != "XXX" {
		t.Errorf("expected XXX for API_KEY, got %q", out["API_KEY"])
	}
	if out["APP_ENV"] != "staging" {
		t.Errorf("expected staging for APP_ENV, got %q", out["APP_ENV"])
	}
	if out["PORT"] != "3000" {
		t.Errorf("expected default PORT=3000, got %q", out["PORT"])
	}
}
