package pipeline_test

import (
	"errors"
	"testing"

	"github.com/yourorg/envsync/internal/pipeline"
)

func TestRun_EmptyPipeline(t *testing.T) {
	input := pipeline.Env{"KEY": "value"}
	out, err := pipeline.New().Run(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value, got %q", out["KEY"])
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := pipeline.Env{"K": "original"}
	pipeline.New().Add(func(env pipeline.Env) error {
		env["K"] = "modified"
		return nil
	}).Run(input) //nolint:errcheck
	if input["K"] != "original" {
		t.Error("input was mutated")
	}
}

func TestRun_StepError_Halts(t *testing.T) {
	sentinel := errors.New("boom")
	called := false
	_, err := pipeline.New().
		Add(func(_ pipeline.Env) error { return sentinel }).
		Add(func(_ pipeline.Env) error { called = true; return nil }).
		Run(pipeline.Env{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if called {
		t.Error("second step should not have been called")
	}
}

func TestStepUpperCaseKeys(t *testing.T) {
	out, err := pipeline.New().
		Add(pipeline.StepUpperCaseKeys()).
		Run(pipeline.Env{"db_host": "localhost", "PORT": "8080"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
	if _, ok := out["db_host"]; ok {
		t.Error("old lowercase key should be gone")
	}
}

func TestStepStripPrefix(t *testing.T) {
	out, err := pipeline.New().
		Add(pipeline.StepStripPrefix("APP_")).
		Run(pipeline.Env{"APP_PORT": "3000", "HOST": "x"})
	if err != nil {
		t.Fatal(err)
	}
	if out["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %q", out["PORT"])
	}
	if _, ok := out["APP_PORT"]; ok {
		t.Error("APP_PORT should have been stripped")
	}
}

func TestStepRequireKeys_Missing(t *testing.T) {
	_, err := pipeline.New().
		Add(pipeline.StepRequireKeys("MUST_EXIST")).
		Run(pipeline.Env{})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestStepSetDefaults(t *testing.T) {
	out, err := pipeline.New().
		Add(pipeline.StepSetDefaults(pipeline.Env{"LOG_LEVEL": "info", "PORT": "8080"})).
		Run(pipeline.Env{"PORT": "9090"})
	if err != nil {
		t.Fatal(err)
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected default LOG_LEVEL=info, got %q", out["LOG_LEVEL"])
	}
	if out["PORT"] != "9090" {
		t.Error("existing PORT should not be overwritten by default")
	}
}

func TestStepMaskValues(t *testing.T) {
	isSensitive := func(k string) bool { return k == "SECRET" }
	out, err := pipeline.New().
		Add(pipeline.StepMaskValues(isSensitive, "***")).
		Run(pipeline.Env{"SECRET": "topsecret", "HOST": "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	if out["SECRET"] != "***" {
		t.Errorf("expected masked value, got %q", out["SECRET"])
	}
	if out["HOST"] != "localhost" {
		t.Error("non-sensitive key should be unchanged")
	}
}
