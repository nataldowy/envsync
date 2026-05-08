package pipeline_test

import (
	"testing"

	"github.com/yourorg/envsync/internal/pipeline"
)

func TestStepMaskValues_EmptyValueNotMasked(t *testing.T) {
	isSensitive := func(k string) bool { return true }
	out, err := pipeline.New().
		Add(pipeline.StepMaskValues(isSensitive, "***")).
		Run(pipeline.Env{"EMPTY_SECRET": ""})
	if err != nil {
		t.Fatal(err)
	}
	if out["EMPTY_SECRET"] != "" {
		t.Errorf("empty value should not be masked, got %q", out["EMPTY_SECRET"])
	}
}

func TestStepRequireKeys_AllPresent(t *testing.T) {
	_, err := pipeline.New().
		Add(pipeline.StepRequireKeys("A", "B")).
		Run(pipeline.Env{"A": "1", "B": "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestChainedSteps(t *testing.T) {
	out, err := pipeline.New().
		Add(pipeline.StepUpperCaseKeys()).
		Add(pipeline.StepSetDefaults(pipeline.Env{"LOG_LEVEL": "warn"})).
		Add(pipeline.StepRequireKeys("DB_HOST")).
		Run(pipeline.Env{"db_host": "pg", "PORT": "5432"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "pg" {
		t.Errorf("DB_HOST mismatch: %q", out["DB_HOST"])
	}
	if out["LOG_LEVEL"] != "warn" {
		t.Errorf("LOG_LEVEL default not applied: %q", out["LOG_LEVEL"])
	}
}
