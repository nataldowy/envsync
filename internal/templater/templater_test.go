package templater

import (
	"strings"
	"testing"
)

func TestRender_BasicSubstitution(t *testing.T) {
	tmpl := "APP_ENV={{APP_ENV}}\nDATABASE_URL={{DATABASE_URL}}"
	vars := map[string]string{
		"APP_ENV":      "production",
		"DATABASE_URL": "postgres://localhost/mydb",
	}
	out, err := Render(tmpl, vars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got: %s", out)
	}
	if !strings.Contains(out, "DATABASE_URL=postgres://localhost/mydb") {
		t.Errorf("expected DATABASE_URL substituted, got: %s", out)
	}
}

func TestRender_CommentPassthrough(t *testing.T) {
	tmpl := "# This is {{NOT_REPLACED}}\nKEY={{KEY}}"
	vars := map[string]string{"KEY": "val"}
	out, err := Render(tmpl, vars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "# This is {{NOT_REPLACED}}") {
		t.Errorf("comment line should be unchanged, got: %s", out)
	}
}

func TestRender_MissingKey_FailOnMissing(t *testing.T) {
	tmpl := "KEY={{MISSING}}"
	_, err := Render(tmpl, map[string]string{}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing placeholder")
	}
	if !strings.Contains(err.Error(), "MISSING") {
		t.Errorf("error should mention key name, got: %v", err)
	}
}

func TestRender_MissingKey_FillEmpty(t *testing.T) {
	tmpl := "KEY={{MISSING}}"
	opts := Options{FailOnMissing: false, FillEmpty: true}
	out, err := Render(tmpl, map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "KEY=" {
		t.Errorf("expected 'KEY=' got %q", out)
	}
}

func TestRender_MissingKey_LeaveToken(t *testing.T) {
	tmpl := "KEY={{MISSING}}"
	opts := Options{FailOnMissing: false, FillEmpty: false}
	out, err := Render(tmpl, map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "KEY={{MISSING}}" {
		t.Errorf("expected token preserved, got %q", out)
	}
}

func TestPlaceholders_ReturnsSorted(t *testing.T) {
	tmpl := "{{ZEBRA}}\n{{ALPHA}}\n{{MIDDLE}}\n{{ALPHA}}"
	got := Placeholders(tmpl)
	want := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	if len(got) != len(want) {
		t.Fatalf("expected %v got %v", want, got)
	}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("index %d: expected %q got %q", i, k, got[i])
		}
	}
}

func TestPlaceholders_NoPlaceholders(t *testing.T) {
	got := Placeholders("KEY=value\n# comment")
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
