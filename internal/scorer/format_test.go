package scorer

import (
	"strings"
	"testing"
)

func TestFormat_NoIssues(t *testing.T) {
	r := Result{Score: 100, Grade: "A", Issues: nil}
	out := Format(r)
	if !strings.Contains(out, "No issues found.") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestFormat_IssuesSortedBySeverity(t *testing.T) {
	r := Result{
		Score: 60,
		Grade: "C",
		Issues: []Issue{
			{Key: "Z_KEY", Severity: "low", Message: "low issue"},
			{Key: "A_KEY", Severity: "high", Message: "high issue"},
			{Key: "M_KEY", Severity: "medium", Message: "medium issue"},
		},
	}
	out := Format(r)
	hiPos := strings.Index(out, "HIGH")
	medPos := strings.Index(out, "MEDIUM")
	lowPos := strings.Index(out, "LOW")
	if hiPos > medPos || medPos > lowPos {
		t.Errorf("issues not sorted by severity: hi=%d med=%d low=%d", hiPos, medPos, lowPos)
	}
}

func TestFormat_IssueCount(t *testing.T) {
	r := Result{
		Score: 70,
		Grade: "B",
		Issues: []Issue{
			{Key: "K1", Severity: "low", Message: "m1"},
			{Key: "K2", Severity: "low", Message: "m2"},
		},
	}
	out := Format(r)
	if !strings.Contains(out, "Issues (2)") {
		t.Errorf("expected 'Issues (2)' in output, got: %s", out)
	}
}

func TestSummary_ZeroScore(t *testing.T) {
	r := Result{Score: 0, Grade: "F", Issues: []Issue{}}
	s := Summary(r)
	if !strings.Contains(s, "grade=F") {
		t.Errorf("expected grade=F in summary, got: %s", s)
	}
}
