package checker

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_DefaultsAndCustom(t *testing.T) {
	dir := t.TempDir()

	defaultsPath := filepath.Join(dir, "defaults.yaml")
	if err := os.WriteFile(defaultsPath, []byte("targets:\n  - example.com:443\n"), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadConfig(defaultsPath)
	if err != nil {
		t.Fatalf("unexpected LoadConfig error: %v", err)
	}
	if len(cfg.Targets) != 1 || cfg.Targets[0] != "example.com:443" {
		t.Fatalf("unexpected targets: %v", cfg.Targets)
	}
	if cfg.CheckInterval != 5*time.Minute {
		t.Fatalf("expected default check interval 5m, got %v", cfg.CheckInterval)
	}
	if cfg.WarnThresholdDays != 14.0 {
		t.Fatalf("expected default warn threshold 14.0, got %v", cfg.WarnThresholdDays)
	}

	customPath := filepath.Join(dir, "custom.yaml")
	customYAML := "targets:\n  - api.example.com:8443\ncheck_interval: 30s\nwarn_threshold_days: 7\nwebhook_url: https://hooks.example.com/alert\n"
	if err := os.WriteFile(customPath, []byte(customYAML), 0o600); err != nil {
		t.Fatalf("failed to write custom config: %v", err)
	}

	custom, err := LoadConfig(customPath)
	if err != nil {
		t.Fatalf("unexpected LoadConfig error: %v", err)
	}
	if custom.CheckInterval != 30*time.Second {
		t.Fatalf("expected custom interval 30s, got %v", custom.CheckInterval)
	}
	if custom.WarnThresholdDays != 7.0 {
		t.Fatalf("expected custom warn threshold 7.0, got %v", custom.WarnThresholdDays)
	}
	if custom.WebhookURL != "https://hooks.example.com/alert" {
		t.Fatalf("unexpected webhook URL: %q", custom.WebhookURL)
	}
}

func TestLoadConfig_Errors(t *testing.T) {
	if _, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Fatal("expected error for missing config file")
	}

	badPath := filepath.Join(t.TempDir(), "bad.yaml")
	if err := os.WriteFile(badPath, []byte("targets: [unclosed"), 0o600); err != nil {
		t.Fatalf("failed to write bad yaml: %v", err)
	}
	if _, err := LoadConfig(badPath); err == nil {
		t.Fatal("expected error for invalid yaml")
	}
}
