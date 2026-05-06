package config

import "testing"

func TestLoadLogFormatDefaultText(t *testing.T) {
	t.Setenv("LOG_FORMAT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.LogFormat != "text" {
		t.Fatalf("expected log format text, got %q", cfg.LogFormat)
	}
}

func TestLoadLogFormatJSON(t *testing.T) {
	t.Setenv("LOG_FORMAT", "json")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.LogFormat != "json" {
		t.Fatalf("expected log format json, got %q", cfg.LogFormat)
	}
}

func TestLoadLogFormatInvalid(t *testing.T) {
	t.Setenv("LOG_FORMAT", "xml")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid LOG_FORMAT")
	}
	if err.Error() != "LOG_FORMAT must be one of: text, json" {
		t.Fatalf("unexpected error: %v", err)
	}
}
