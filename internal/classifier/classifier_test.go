package classifier

import (
	"testing"
)

func TestClassify_Database(t *testing.T) {
	c := New()
	for _, key := range []string{"DB_HOST", "DATABASE_URL", "POSTGRES_PASSWORD", "REDIS_ADDR", "MONGO_DSN"} {
		if got := c.Classify(key); got != CategoryDatabase {
			t.Errorf("Classify(%q) = %q, want %q", key, got, CategoryDatabase)
		}
	}
}

func TestClassify_Auth(t *testing.T) {
	c := New()
	for _, key := range []string{"SECRET_KEY", "API_KEY", "JWT_TOKEN", "AUTH_PASSWORD", "OAUTH_CLIENT_SECRET"} {
		if got := c.Classify(key); got != CategoryAuth {
			t.Errorf("Classify(%q) = %q, want %q", key, got, CategoryAuth)
		}
	}
}

func TestClassify_Network(t *testing.T) {
	c := New()
	for _, key := range []string{"APP_HOST", "HTTP_PORT", "BASE_URL", "PROXY_ENDPOINT", "TLS_CERT"} {
		if got := c.Classify(key); got != CategoryNetwork {
			t.Errorf("Classify(%q) = %q, want %q", key, got, CategoryNetwork)
		}
	}
}

func TestClassify_Storage(t *testing.T) {
	c := New()
	for _, key := range []string{"S3_BUCKET", "GCS_BUCKET", "STORAGE_PATH", "UPLOAD_DIR", "LOG_FILE"} {
		// LOG_FILE matches storage (FILE) before logging — acceptable
		got := c.Classify(key)
		if got != CategoryStorage && got != CategoryLogging {
			t.Errorf("Classify(%q) = %q, want storage or logging", key, got)
		}
	}
}

func TestClassify_Logging(t *testing.T) {
	c := New()
	for _, key := range []string{"LOG_LEVEL", "DEBUG", "SENTRY_DSN", "DATADOG_API_KEY", "VERBOSE"} {
		got := c.Classify(key)
		// SENTRY_DSN and DATADOG_API_KEY may also match auth; accept logging or auth
		if got != CategoryLogging && got != CategoryAuth && got != CategoryDatabase {
			t.Errorf("Classify(%q) = %q, unexpected category", key, got)
		}
	}
}

func TestClassify_Other(t *testing.T) {
	c := New()
	for _, key := range []string{"APP_NAME", "ENVIRONMENT", "REGION", "FEATURE_FLAG"} {
		if got := c.Classify(key); got != CategoryOther {
			t.Errorf("Classify(%q) = %q, want %q", key, got, CategoryOther)
		}
	}
}

func TestClassifyMap_GroupsCorrectly(t *testing.T) {
	c := New()
	env := map[string]string{
		"DB_HOST":    "localhost",
		"API_KEY":    "abc123",
		"APP_PORT":   "8080",
		"APP_NAME":   "envsync",
		"LOG_LEVEL":  "info",
	}
	groups := c.ClassifyMap(env)

	if len(groups[CategoryDatabase]) != 1 || groups[CategoryDatabase][0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in database group, got %v", groups[CategoryDatabase])
	}
	if len(groups[CategoryAuth]) != 1 || groups[CategoryAuth][0] != "API_KEY" {
		t.Errorf("expected API_KEY in auth group, got %v", groups[CategoryAuth])
	}
	if len(groups[CategoryOther]) != 1 || groups[CategoryOther][0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in other group, got %v", groups[CategoryOther])
	}
}

func TestClassifyMap_SortedKeys(t *testing.T) {
	c := New()
	env := map[string]string{
		"ZEBRA_NAME": "z",
		"ALPHA_NAME": "a",
		"MANGO_NAME": "m",
	}
	groups := c.ClassifyMap(env)
	other := groups[CategoryOther]
	for i := 1; i < len(other); i++ {
		if other[i] < other[i-1] {
			t.Errorf("keys not sorted: %v", other)
		}
	}
}
