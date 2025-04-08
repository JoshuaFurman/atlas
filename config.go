package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ModelConfig represents the configuration for a specific language model
type ModelConfig struct {
	Name         string  `yaml:"name"`
	Temperature  float32 `yaml:"temp"`
	SystemPrompt string  `yaml:"system_prompt"`
}

// ProviderConfig represents the configuration for an AI provider
type ProviderConfig struct {
	Endpoint string        `yaml:"endpoint"`
	APIKey   string        `yaml:"api_key"`
	Models   []ModelConfig `yaml:"models"`
}

// Config represents the root configuration structure with dynamic provider names
type Config struct {
	ActiveProvider string
	ActiveModel    string
	Providers      map[string]ProviderConfig `yaml:",inline"`
}

// LoadConfig loads the configuration from the default path
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".config", "atlas", "config.yml")
	return LoadConfigFromPath(configPath)
}

// LoadConfigFromPath loads the configuration from a specific path
func LoadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var config Config
	// Initialize the providers map
	config.Providers = make(map[string]ProviderConfig)

	if err := yaml.Unmarshal(data, &config.Providers); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// GetAllProviders returns a list of all provider names
func (c *Config) GetAllProviders() []string {
	providers := make([]string, 0, len(c.Providers))
	for name := range c.Providers {
		providers = append(providers, name)
	}
	return providers
}

// GetModelConfig retrieves a model config by provider and model name
func (c *Config) GetModelConfig(provider, modelName string) (*ModelConfig, error) {
	providerConfig, exists := c.Providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}

	for _, model := range providerConfig.Models {
		if model.Name == modelName {
			return &model, nil
		}
	}

	return nil, fmt.Errorf("model %s not found for provider %s", modelName, provider)
}

// GetProviderConfig retrieves a provider's configuration by name
func (c *Config) GetProviderConfig(provider string) (*ProviderConfig, error) {
	providerConfig, exists := c.Providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}
	return &providerConfig, nil
}

// GetModelsForProvider returns all models for a given provider
func (c *Config) GetModelsForProvider(provider string) ([]ModelConfig, error) {
	providerConfig, exists := c.Providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", provider)
	}
	return providerConfig.Models, nil
}
