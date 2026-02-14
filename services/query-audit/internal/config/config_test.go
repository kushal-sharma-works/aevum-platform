package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	_ = os.Unsetenv("SERVER_PORT")
	_ = os.Unsetenv("SERVER_HOST")
	_ = os.Unsetenv("ELASTICSEARCH_URLS")
	_ = os.Unsetenv("EVENT_TIMELINE_URL")
	_ = os.Unsetenv("DECISION_ENGINE_URL")
	_ = os.Unsetenv("SYNC_INTERVAL")
	_ = os.Unsetenv("SYNC_MAX_BACKOFF")
	_ = os.Unsetenv("SYNC_BATCH_SIZE")
	_ = os.Unsetenv("ENVIRONMENT")

	cfg := Load()
	if cfg.Server.Port != 8080 {
		t.Fatalf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Fatalf("expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}
	if len(cfg.Elasticsearch.URLs) != 1 || cfg.Elasticsearch.URLs[0] != "http://localhost:9200" {
		t.Fatalf("unexpected default elasticsearch urls: %+v", cfg.Elasticsearch.URLs)
	}
	if cfg.Sync.Interval != 30*time.Second {
		t.Fatalf("expected default sync interval 30s, got %s", cfg.Sync.Interval)
	}
	if cfg.Sync.MaxBackoff != 300*time.Second {
		t.Fatalf("expected default sync max backoff 300s, got %s", cfg.Sync.MaxBackoff)
	}
	if cfg.Sync.BatchSize != 500 {
		t.Fatalf("expected default batch size 500, got %d", cfg.Sync.BatchSize)
	}
	if cfg.Environment != "development" {
		t.Fatalf("expected default environment development, got %s", cfg.Environment)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("SERVER_PORT", "9099")
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("ELASTICSEARCH_URLS", "http://es1:9200,http://es2:9200")
	t.Setenv("EVENT_TIMELINE_URL", "http://event-timeline:8080")
	t.Setenv("DECISION_ENGINE_URL", "http://decision-engine:8081")
	t.Setenv("SYNC_INTERVAL", "5")
	t.Setenv("SYNC_MAX_BACKOFF", "60")
	t.Setenv("SYNC_BATCH_SIZE", "100")
	t.Setenv("ENVIRONMENT", "sit")

	cfg := Load()
	if cfg.Server.Port != 9099 {
		t.Fatalf("expected port 9099, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected host 127.0.0.1, got %s", cfg.Server.Host)
	}
	if len(cfg.Elasticsearch.URLs) != 2 {
		t.Fatalf("expected 2 elasticsearch URLs, got %d", len(cfg.Elasticsearch.URLs))
	}
	if cfg.Sync.Interval != 5*time.Second {
		t.Fatalf("expected sync interval 5s, got %s", cfg.Sync.Interval)
	}
	if cfg.Sync.MaxBackoff != 60*time.Second {
		t.Fatalf("expected sync max backoff 60s, got %s", cfg.Sync.MaxBackoff)
	}
	if cfg.Sync.BatchSize != 100 {
		t.Fatalf("expected batch size 100, got %d", cfg.Sync.BatchSize)
	}
	if cfg.Environment != "sit" {
		t.Fatalf("expected environment sit, got %s", cfg.Environment)
	}
}

func TestLoadFallsBackOnInvalidNumbers(t *testing.T) {
	t.Setenv("SERVER_PORT", "not-a-number")
	t.Setenv("SYNC_INTERVAL", "oops")
	t.Setenv("SYNC_MAX_BACKOFF", "invalid")
	t.Setenv("SYNC_BATCH_SIZE", "broken")

	cfg := Load()
	if cfg.Server.Port != 8080 {
		t.Fatalf("expected fallback port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Sync.Interval != 30*time.Second {
		t.Fatalf("expected fallback interval 30s, got %s", cfg.Sync.Interval)
	}
	if cfg.Sync.MaxBackoff != 300*time.Second {
		t.Fatalf("expected fallback max backoff 300s, got %s", cfg.Sync.MaxBackoff)
	}
	if cfg.Sync.BatchSize != 500 {
		t.Fatalf("expected fallback batch size 500, got %d", cfg.Sync.BatchSize)
	}
}
