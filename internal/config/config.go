package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds user configuration
type Config struct {
	DefaultConsultant string  `json:"default_consultant"`
	DefaultClient     string  `json:"default_client"`
	DefaultProject    string  `json:"default_project"`
	DefaultRate       float64 `json:"default_rate"`
	Database          Database `json:"database"`
}

// Database holds database configuration
type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".worklog")
	configFile := filepath.Join(configDir, "config.json")
	return configFile, nil
}

// Load reads the config from file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return &Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Save writes the config to file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// SetDefaults updates the default values
func (c *Config) SetDefaults(consultant, client, project string, rate float64) {
	if consultant != "" {
		c.DefaultConsultant = consultant
	}
	if client != "" {
		c.DefaultClient = client
	}
	if project != "" {
		c.DefaultProject = project
	}
	if rate > 0 {
		c.DefaultRate = rate
	}
}

// ClearDefaults clears the saved configuration
func (c *Config) ClearDefaults() {
	c.DefaultConsultant = ""
	c.DefaultClient = ""
	c.DefaultProject = ""
	c.DefaultRate = 0
	c.Database = Database{}
}
