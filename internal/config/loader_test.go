package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	for _, key := range []string{
		"APP_CLUSTER",
		"APP_ENVIRONMENT",
		"APP_PORT",
		"LOGGER_LEVEL",
		"LOGGER_ENABLE",
		"SENTRY_ENABLE",
		"SENTRY_DSN",
		"SENTRY_DEBUG",
	} {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset env var %s: %v", key, err)
		}
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Cluster != "dev" {
		t.Errorf("expected default cluster 'dev', got %s", cfg.Cluster)
	}
	if cfg.Environment != "local" {
		t.Errorf("expected default environment 'local', got %s", cfg.Environment)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.Logger.Level != "INFO" {
		t.Errorf("expected default logger level INFO, got %s", cfg.Logger.Level)
	}
	if cfg.Logger.Enable != false {
		t.Errorf("expected default logger enable false, got %v", cfg.Logger.Enable)
	}
	if cfg.Sentry.Enable != false {
		t.Errorf("expected default sentry enable false, got %v", cfg.Sentry.Enable)
	}
	if cfg.Sentry.Dsn != "" {
		t.Errorf("expected default sentry dsn empty, got %s", cfg.Sentry.Dsn)
	}
	if cfg.Sentry.Debug != false {
		t.Errorf("expected default sentry debug false, got %v", cfg.Sentry.Debug)
	}
}

func TestLoadWithEnv(t *testing.T) {
	t.Setenv("APP_CLUSTER", "prod")
	t.Setenv("APP_ENVIRONMENT", "staging")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("LOGGER_LEVEL", "DEBUG")
	t.Setenv("LOGGER_ENABLE", "true")
	t.Setenv("SENTRY_ENABLE", "true")
	t.Setenv("SENTRY_DSN", "dsn")
	t.Setenv("SENTRY_DEBUG", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Cluster != "prod" {
		t.Errorf("expected cluster 'prod', got %s", cfg.Cluster)
	}
	if cfg.Environment != "staging" {
		t.Errorf("expected environment 'staging', got %s", cfg.Environment)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.Logger.Level != "DEBUG" {
		t.Errorf("expected logger level DEBUG, got %s", cfg.Logger.Level)
	}
	if cfg.Logger.Enable != true {
		t.Errorf("expected logger enable true, got %v", cfg.Logger.Enable)
	}
	if cfg.Sentry.Enable != true {
		t.Errorf("expected sentry enable true, got %v", cfg.Sentry.Enable)
	}
	if cfg.Sentry.Dsn != "dsn" {
		t.Errorf("expected sentry dsn 'dsn', got %s", cfg.Sentry.Dsn)
	}
	if cfg.Sentry.Debug != true {
		t.Errorf("expected sentry debug true, got %v", cfg.Sentry.Debug)
	}
}
