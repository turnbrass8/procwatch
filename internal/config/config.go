package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level procwatch configuration.
type Config struct {
	WebhookURL string          `yaml:"webhook_url"`
	Interval   time.Duration   `yaml:"interval"`
	Processes  []ProcessConfig `yaml:"processes"`
}

// ProcessConfig defines a single process to monitor.
type ProcessConfig struct {
	Name       string  `yaml:"name"`
	Match      string  `yaml:"match"`       // substring or exact process name to match
	MaxCPU     float64 `yaml:"max_cpu"`     // percent, 0 = disabled
	MaxMemMB   float64 `yaml:"max_mem_mb"`  // megabytes, 0 = disabled
	AlertOnExit bool   `yaml:"alert_on_exit"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.WebhookURL == "" {
		return fmt.Errorf("config: webhook_url is required")
	}
	if c.Interval <= 0 {
		c.Interval = 15 * time.Second
	}
	if len(c.Processes) == 0 {
		return fmt.Errorf("config: at least one process must be defined")
	}
	for i, p := range c.Processes {
		if p.Name == "" {
			return fmt.Errorf("config: process[%d]: name is required", i)
		}
		if p.Match == "" {
			return fmt.Errorf("config: process[%d] %q: match is required", i, p.Name)
		}
	}
	return nil
}
