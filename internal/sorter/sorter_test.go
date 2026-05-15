package sorter_test

import (
	"testing"

	"github.com/user/envsync/internal/sorter"
)

func makeEnv(keys ...string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "value"
	}
	return m
}

func TestSort_Alpha(t *testing.T) {
	env := makeEnv("ZEBRA", "APPLE", "MANGO")
	keys := sorter.Sort(env, sorter.DefaultOptions())
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Fatalf("unexpected alpha order: %v", keys)
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	env := makeEnv("ZEBRA", "APPLE", "MANGO")
	keys := sorter.Sort(env, sorter.Options{Order: sorter.OrderAlphaDesc})
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Fatalf("unexpected alpha_desc order: %v", keys)
	}
}

func TestSort_Length(t *testing.T) {
	env := makeEnv("AB", "ABCDE", "ABC")
	keys := sorter.Sort(env, sorter.Options{Order: sorter.OrderLength})
	if keys[0] != "AB" || keys[1] != "ABC" || keys[2] != "ABCDE" {
		t.Fatalf("unexpected length order: %v", keys)
	}
}

func TestSort_Grouped_SensitiveLast(t *testing.T) {
	env := makeEnv("DB_HOST", "API_SECRET", "APP_NAME", "DB_PASSWORD")
	keys := sorter.Sort(env, sorter.Options{Order: sorter.OrderGrouped})

	sensitiveIdx := -1
	plainLastIdx := -1
	for i, k := range keys {
		switch k {
		case "API_SECRET", "DB_PASSWORD":
			if sensitiveIdx == -1 {
				sensitiveIdx = i
			}
		default:
			plainLastIdx = i
		}
	}
	if plainLastIdx > sensitiveIdx && sensitiveIdx != -1 {
		t.Fatalf("plain key appeared after sensitive key: %v", keys)
	}
}

func TestSort_Grouped_SensitiveFirst(t *testing.T) {
	env := makeEnv("DB_HOST", "API_TOKEN", "APP_NAME")
	keys := sorter.Sort(env, sorter.Options{
		Order:          sorter.OrderGrouped,
		SensitiveFirst: true,
	})
	if keys[0] != "API_TOKEN" {
		t.Fatalf("expected sensitive key first, got %v", keys)
	}
}

func TestSort_EmptyMap(t *testing.T) {
	keys := sorter.Sort(map[string]string{}, sorter.DefaultOptions())
	if len(keys) != 0 {
		t.Fatalf("expected empty slice, got %v", keys)
	}
}

func TestSort_StableOnTie(t *testing.T) {
	// Two keys of same length — should fall back to alpha
	env := makeEnv("ZZ", "AA")
	keys := sorter.Sort(env, sorter.Options{Order: sorter.OrderLength})
	if keys[0] != "AA" || keys[1] != "ZZ" {
		t.Fatalf("expected alpha tie-break, got %v", keys)
	}
}
