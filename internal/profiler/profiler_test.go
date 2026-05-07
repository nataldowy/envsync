package profiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze_TotalKeys(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	p := Analyze(env, nil)
	assert.Equal(t, 2, p.TotalKeys)
}

func TestAnalyze_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_TOKEN":   "tok123",
		"APP_NAME":    "myapp",
	}
	p := Analyze(env, nil)
	assert.ElementsMatch(t, []string{"DB_PASSWORD", "API_TOKEN"}, p.SensitiveKeys)
}

func TestAnalyze_EmptyValues(t *testing.T) {
	env := map[string]string{
		"EMPTY_KEY": "",
		"BLANK_KEY": "   ",
		"SET_KEY":   "value",
	}
	p := Analyze(env, nil)
	assert.ElementsMatch(t, []string{"EMPTY_KEY", "BLANK_KEY"}, p.EmptyValues)
}

func TestAnalyze_DuplicateKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	keys := []string{"FOO", "BAR", "FOO"} // FOO appears twice
	p := Analyze(env, keys)
	assert.Equal(t, []string{"FOO"}, p.DuplicateKeys)
}

func TestAnalyze_LongestKey(t *testing.T) {
	env := map[string]string{
		"SHORT":            "a",
		"MUCH_LONGER_KEY":  "b",
		"MID_KEY":          "c",
	}
	p := Analyze(env, nil)
	assert.Equal(t, "MUCH_LONGER_KEY", p.LongestKey)
}

func TestAnalyze_EmptyMap(t *testing.T) {
	p := Analyze(map[string]string{}, nil)
	assert.Equal(t, 0, p.TotalKeys)
	assert.Empty(t, p.SensitiveKeys)
	assert.Empty(t, p.EmptyValues)
	assert.Empty(t, p.DuplicateKeys)
	assert.Equal(t, "", p.LongestKey)
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	assert.True(t, isSensitive("db_Password"))
	assert.True(t, isSensitive("PRIVATE_KEY"))
	assert.False(t, isSensitive("APP_NAME"))
	assert.False(t, isSensitive("PORT"))
}
