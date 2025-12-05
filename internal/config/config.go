package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds user configuration
type Config struct {
	DefaultConsultant string   `mapstructure:"default_consultant"`
	DefaultClient     string   `mapstructure:"default_client"`
	DefaultProject    string   `mapstructure:"default_project"`
	DefaultRate       float64  `mapstructure:"default_rate"`
	Database          Database `mapstructure:"database"`
}

// Database holds database configuration
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

var v *viper.Viper

// Initialize loads configuration from config files
func Initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".worklog")

	v = viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName("config")
	v.SetConfigType("json")

	// Load config
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Bind environment variables - they override config file values
	v.SetEnvPrefix("WORKLOG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind all environment variables to config keys
	v.BindEnv("database.host")
	v.BindEnv("database.port")
	v.BindEnv("database.user")
	v.BindEnv("database.password")
	v.BindEnv("database.name")
	v.BindEnv("default_consultant")
	v.BindEnv("default_client")
	v.BindEnv("default_project")
	v.BindEnv("default_rate")

	return nil
}

// Get returns the current configuration
func Get() (*Config, error) {
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}

// SaveDefaults writes default values to config file
func SaveDefaults(consultant, client, project string, rate float64) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".worklog", "config.json")

	// Update viper instance
	if consultant != "" {
		v.Set("default_consultant", consultant)
	}
	if client != "" {
		v.Set("default_client", client)
	}
	if project != "" {
		v.Set("default_project", project)
	}
	if rate > 0 {
		v.Set("default_rate", rate)
	}

	// Write to file
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// ClearDefaults removes defaults from config file
func ClearDefaults() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".worklog", "config.json")

	v.Set("default_consultant", "")
	v.Set("default_client", "")
	v.Set("default_project", "")
	v.Set("default_rate", 0)

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to clear config: %w", err)
	}

	return nil
}
