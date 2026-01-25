package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents the server configuration
type Config struct {
	values map[string]interface{}
	Mods   ModsConfig `yaml:"mods"`
}

// ModsConfig represents mod configuration
type ModsConfig struct {
	Enabled bool        `yaml:"enabled"`
	Mods    []ModConfig `yaml:"mods"`
}

// ModConfig represents a single mod
type ModConfig struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	URL         string            `yaml:"url"`
	Checksum    string            `yaml:"checksum"`
	InstallPath string            `yaml:"installPath"`
	LoadOrder   int               `yaml:"loadOrder"`
	ConfigFiles []ModConfigFile   `yaml:"configFiles"`
	Metadata    map[string]string `yaml:"metadata"`
}

// ModConfigFile represents a mod configuration file
type ModConfigFile struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Template    bool   `yaml:"template"`
}

// Load reads configuration from a YAML file and environment variables
func Load(path string) (*Config, error) {
	cfg := &Config{
		values: make(map[string]interface{}),
	}

	// Read from file if it exists
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg.values); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Load mods config if separate file exists
	modsPath := "./config-data/mods.yaml"
	if _, err := os.Stat(modsPath); err == nil {
		data, err := os.ReadFile(modsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read mods config: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg.Mods); err != nil {
			return nil, fmt.Errorf("failed to parse mods config: %w", err)
		}
	}

	return cfg, nil
}

// GetString returns a string value from config or environment
func (c *Config) GetString(key string, defaultValue string) string {
	// Check environment first
	if val := os.Getenv(key); val != "" {
		return val
	}

	// Check config file
	if val, ok := c.values[key].(string); ok {
		return val
	}

	return defaultValue
}

// GetBool returns a boolean value from config or environment
func (c *Config) GetBool(key string, defaultValue bool) bool {
	// Check environment first
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
		return val == "TRUE" || val == "true" || val == "1"
	}

	// Check config file
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			return v == "TRUE" || v == "true" || v == "1"
		}
	}

	return defaultValue
}

// GetInt returns an integer value from config or environment
func (c *Config) GetInt(key string, defaultValue int) int {
	// Check environment first
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}

	// Check config file
	if val, ok := c.values[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}

	return defaultValue
}

// Get returns a raw value
func (c *Config) Get(key string) interface{} {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return c.values[key]
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.values[key] = value
}
