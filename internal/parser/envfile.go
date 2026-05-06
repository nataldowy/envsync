package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline or preceding comment
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
}

// Parse reads a .env file and returns an EnvFile.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	ef := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	var pendingComment string

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			pendingComment = strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			continue
		}

		key, value, found := strings.Cut(trimmed, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = stripInlineComment(value)
		value = unquote(value)

		ef.Entries = append(ef.Entries, Entry{
			Key:     key,
			Value:   value,
			Comment: pendingComment,
			Line:    lineNum,
		})
		pendingComment = ""
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return ef, nil
}

// ToMap converts entries to a key→value map.
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		m[e.Key] = e.Value
	}
	return m
}

func stripInlineComment(s string) string {
	if idx := strings.Index(s, " #"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
