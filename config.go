package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Model represents an AI model configuration
type Model struct {
	Name        string  `yaml:"name"`
	Temperature float64 `yaml:"temperature"`
}

// Provider represents a provider configuration
type Provider struct {
	Endpoint string  `yaml:"endpoint"`
	APIKey   string  `yaml:"api_key"`
	Models   []Model `yaml:"models"`
}

// Config represents the application configuration
type Config struct {
	Providers map[string]Provider `yaml:",inline"`
}

// LoadConfig reads configuration from the specified file path
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// DefaultConfigPath returns the default path to the config file
func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "atlas")
	configPath := filepath.Join(configDir, "config.yml")

	return configPath, nil
}

// LoadDefaultConfig attempts to load the configuration from the default path
func LoadDefaultConfig() (*Config, error) {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}

	return LoadConfig(configPath)
}

// EnsureConfigExists checks if the config file exists, and creates a default one if it doesn't
func EnsureConfigExists() error {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	// Check if config file exists
	_, err = os.Stat(configPath)
	if err == nil {
		// File exists, nothing to do
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error checking config file: %w", err)
	}

	// Config directory
	configDir := filepath.Dir(configPath)

	// Create directory if it doesn't exist
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a default config
	defaultConfig := Config{
		Providers: map[string]Provider{
			"example-provider": {
				Endpoint: "https://api.example.com/v1",
				APIKey:   "YOUR_API_KEY_HERE",
				Models: []Model{
					{
						Name:        "default-model",
						Temperature: 0.7,
					},
				},
			},
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("error creating default config: %w", err)
	}

	// Write the default config file
	err = os.WriteFile(configPath, data, 0600) // Secure permissions since it may contain API keys
	if err != nil {
		return fmt.Errorf("error writing default config file: %w", err)
	}

	return nil
}

// GetProvider returns a provider configuration by name
func (c *Config) GetProvider(name string) (Provider, error) {
	provider, exists := c.Providers[name]
	if !exists {
		return Provider{}, fmt.Errorf("provider '%s' not found in configuration", name)
	}
	return provider, nil
}

// SaveConfig saves the current configuration to the specified file path
func (c *Config) SaveConfig(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600) // Secure permissions for API keys
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// SaveDefaultConfig saves the configuration to the default path
func (c *Config) SaveDefaultConfig() error {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	return c.SaveConfig(configPath)
}

// AddProvider adds or updates a provider in the configuration
func (c *Config) AddProvider(name string, provider Provider) {
	if c.Providers == nil {
		c.Providers = make(map[string]Provider)
	}
	c.Providers[name] = provider
}
