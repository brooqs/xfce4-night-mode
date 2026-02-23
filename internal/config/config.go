package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	appName    = "xfce4-night-mode"
	configFile = "config.yaml"
)

// ThemeConfig holds the theme names for a mode (day or night).
type ThemeConfig struct {
	GtkTheme  string `yaml:"gtk_theme"`
	IconTheme string `yaml:"icon_theme"`
	WmTheme   string `yaml:"wm_theme"`
}

// Location holds GPS coordinates.
type Location struct {
	Latitude  float64 `yaml:"latitude"`
	Longitude float64 `yaml:"longitude"`
}

// Config is the root configuration structure.
type Config struct {
	Location      Location    `yaml:"location"`
	DayTheme      ThemeConfig `yaml:"day_theme"`
	NightTheme    ThemeConfig `yaml:"night_theme"`
	CheckInterval int         `yaml:"check_interval"` // minutes
}

// DefaultConfig returns a sensible default configuration (Istanbul).
func DefaultConfig() *Config {
	return &Config{
		Location: Location{
			Latitude:  41.0082,
			Longitude: 28.9784,
		},
		DayTheme: ThemeConfig{
			GtkTheme:  "Adwaita",
			IconTheme: "Adwaita",
			WmTheme:   "Default",
		},
		NightTheme: ThemeConfig{
			GtkTheme:  "Adwaita-dark",
			IconTheme: "Adwaita",
			WmTheme:   "Default-hdpi",
		},
		CheckInterval: 5,
	}
}

// ConfigDir returns the path to the configuration directory.
func ConfigDir() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, appName), nil
}

// ConfigPath returns the full path to the config file.
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

// Load reads the config file from the default location.
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFrom(path)
}

// LoadFrom reads the config file from a specific path.
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	if cfg.CheckInterval < 1 {
		cfg.CheckInterval = 1
	}

	return cfg, nil
}

// Init creates the default config file if it doesn't already exist.
func Init() (string, error) {
	path, err := ConfigPath()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		return path, fmt.Errorf("config file already exists: %s", path)
	}

	data, err := yaml.Marshal(DefaultConfig())
	if err != nil {
		return "", fmt.Errorf("could not marshal default config: %w", err)
	}

	header := "# xfce4-night-mode configuration\n# Edit this file to set your location and preferred themes.\n\n"
	if err := os.WriteFile(path, []byte(header+string(data)), 0o644); err != nil {
		return "", fmt.Errorf("could not write config file: %w", err)
	}

	return path, nil
}
