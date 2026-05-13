package scorer

import (
	"strings"
	"testing"
)

func TestScore_PerfectEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	r := Score(env)
	if r.Score != 100 {
		t.Errorf("expected 100, got %d", r.Score)
	}
	if r.Grade != "A" {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
	if len(r.Issues) != 0 {
		t.Errorf("expected no issues, got %v", r.Issues)
	}
}

func TestScore_SensitiveEmptyValue(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": ""}
	r := Score(env)
	if r.Score >= 100 {
		t.Error("expected penalty for empty sensitive value")
	}
	found := false
	for _, iss := range r.Issues {
		if iss.Key == "DB_PASSWORD" && iss.Severity == "high" {
			found = true
		}
	}
	if !found {
		t.Error("expected high-severity issue for DB_PASSWORD")
	}
}

func TestScore_LowercaseKey(t *testing.T) {
	env := map[string]string{"app_name": "x"}
	r := Score(env)
	if r.Score == 100 {
		t.Error("expected penalty for lowercase key")
	}
	for _, iss := range r.Issues {
		if iss.Key == "app_name" && iss.Severity == "low" {
			return
		}
	}
	t.Error("expected low-severity issue for app_name")
}

func TestScore_ShortSensitiveValue(t *testing.T) {
	env := map[string]string{"API_KEY": "abc"}
	r := Score(env)
	for _, iss := range r.Issues {
		if iss.Key == "API_KEY" && iss.Severity == "high" {
			return
		}
	}
	t.Error("expected high issue for short API_KEY value")
}

func TestScore_GradeThresholds(t *testing.T) {
	for _, tc := range []struct {
		score int
		want  string
	}{
		{95, "A"}, {80, "B"}, {65, "C"}, {50, "D"}, {30, "F"},
	} {
		g := grade(tc.score)
		if g != tc.want {
			t.Errorf("grade(%d) = %s, want %s", tc.score, g, tc.want)
		}
	}
}

func TestFormat_ContainsGrade(t *testing.T) {
	env := map[string]string{"SECRET": ""}
	r := Score(env)
	out := Format(r)
	if !strings.Contains(out, "Grade:") {
		t.Error("Format output missing Grade label")
	}
	if !strings.Contains(out, "HIGH") {
		t.Error("Format output missing HIGH severity label")
	}
}

func TestSummary_Format(t *testing.T) {
	r := Result{Score: 72, Grade: "B", Issues: []Issue{{Key: "X", Severity: "low", Message: "test"}}}
	s := Summary(r)
	if !strings.Contains(s, "grade=B") || !strings.Contains(s, "score=72") || !strings.Contains(s, "issues=1") {
		t.Errorf("unexpected summary: %s", s)
	}
}
