package checker

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Targets           []string      `yaml:"targets"`
	CheckInterval     time.Duration `yaml:"check_interval"`
	WarnThresholdDays float64       `yaml:"warn_threshold_days"`
	WebhookURL        string        `yaml:"webhook_url"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = 5 * time.Minute
	}
	if cfg.WarnThresholdDays == 0 {
		cfg.WarnThresholdDays = 14.0
	}

	return &cfg, nil
}
