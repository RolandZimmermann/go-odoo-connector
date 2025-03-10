package odoo

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the Odoo connection configuration
type Config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
	DB       string `json:"db"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.URL == "" {
		return nil, fmt.Errorf("URL is required in config")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("username is required in config")
	}
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required in config")
	}
	if config.DB == "" {
		return nil, fmt.Errorf("database name is required in config")
	}

	return &config, nil
}

// NewConnectorFromConfig creates a new Odoo connector using configuration
func NewConnectorFromConfig(configPath string) (*Connector, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return NewConnector(config.URL, config.Username, config.APIKey, config.DB)
}
