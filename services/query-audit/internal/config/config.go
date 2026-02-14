package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server         ServerConfig
	Elasticsearch  ElasticsearchConfig
	EventTimeline  EventTimelineConfig
	DecisionEngine DecisionEngineConfig
	Sync           SyncConfig
	Environment    string
}

// ServerConfig represents server settings
type ServerConfig struct {
	Port int
	Host string
}

// ElasticsearchConfig represents ES settings
type ElasticsearchConfig struct {
	URLs []string
}

// EventTimelineConfig represents Event Timeline Service settings
type EventTimelineConfig struct {
	BaseURL string
}

// DecisionEngineConfig represents Decision Engine settings
type DecisionEngineConfig struct {
	BaseURL string
}

// SyncConfig represents sync settings
type SyncConfig struct {
	Interval   time.Duration
	MaxBackoff time.Duration
	BatchSize  int
}

// Load loads configuration from environment
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Elasticsearch: ElasticsearchConfig{
			URLs: getEnvSlice("ELASTICSEARCH_URLS", []string{"http://localhost:9200"}),
		},
		EventTimeline: EventTimelineConfig{
			BaseURL: getEnv("EVENT_TIMELINE_URL", "http://localhost:8081"),
		},
		DecisionEngine: DecisionEngineConfig{
			BaseURL: getEnv("DECISION_ENGINE_URL", "http://localhost:8082"),
		},
		Sync: SyncConfig{
			Interval:   time.Duration(getEnvInt("SYNC_INTERVAL", 30)) * time.Second,
			MaxBackoff: time.Duration(getEnvInt("SYNC_MAX_BACKOFF", 300)) * time.Second,
			BatchSize:  getEnvInt("SYNC_BATCH_SIZE", 500),
		},
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv gets an environment variable with a default
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

// getEnvInt gets an integer environment variable with a default
func getEnvInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvSlice gets a comma-separated environment variable as a slice
func getEnvSlice(key string, defaultVal []string) []string {
	if val, exists := os.LookupEnv(key); exists {
		return strings.Split(val, ",")
	}
	return defaultVal
}
