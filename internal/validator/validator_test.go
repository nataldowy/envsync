package validator

import (
	"testing"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	issues := Validate(env, nil, Rules{})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidate_InvalidKeyName(t *testing.T) {
	env := map[string]string{"123INVALID": "value"}
	issues := Validate(env, nil, Rules{})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "123INVALID" {
		t.Errorf("unexpected key in issue: %s", issues[0].Key)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{"SECRET": ""}
	issues := Validate(env, nil, Rules{NoEmptyValues: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Message != "empty value" {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestValidate_EmptyValueAllowed(t *testing.T) {
	env := map[string]string{"OPTIONAL": ""}
	issues := Validate(env, nil, Rules{NoEmptyValues: false})
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidate_RequiredKeysMissing(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	issues := Validate(env, nil, Rules{RequiredKeys: []string{"PORT", "DATABASE_URL"}})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected missing key: %s", issues[0].Key)
	}
}

func TestValidate_LineNumbers(t *testing.T) {
	lines := []string{
		"# comment",
		"APP_NAME=myapp",
		"EMPTY_VAL=",
	}
	env := map[string]string{"APP_NAME": "myapp", "EMPTY_VAL": ""}
	issues := Validate(env, lines, Rules{NoEmptyValues: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Line != 3 {
		t.Errorf("expected line 3, got %d", issues[0].Line)
	}
}

func TestIssue_String_WithLine(t *testing.T) {
	i := Issue{Line: 5, Key: "FOO", Message: "empty value"}
	got := i.String()
	want := "line 5: FOO: empty value"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIssue_String_NoLine(t *testing.T) {
	i := Issue{Key: "BAR", Message: "required key is missing"}
	got := i.String()
	want := "BAR: required key is missing"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
