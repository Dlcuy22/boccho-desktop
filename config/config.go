package config

/*
config.go - Configuration management for boccho-ui

Functions:
- GetDefaultConfig: Returns default configuration with standard paths
- LoadConfig: Loads config from boccho.config.json or creates default
- SaveConfig: Saves current config to boccho.config.json
- GetConfigPath: Returns the path to boccho.config.json
- getDefaultFramesPath: Returns default frames path
- GetAppDataDir: Returns app data directory for current OS
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	FramesPath string `json:"framesPath"`
}

func GetAppDataDir() string {
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			homeDir, _ := os.UserHomeDir()
			localAppData = filepath.Join(homeDir, "AppData", "Local")
		}
		return filepath.Join(localAppData, "boccho-ui")
	}

	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".boccho-ui")
}

func getDefaultFramesPath() string {
	return filepath.Join(GetAppDataDir(), "Frames")
}

func GetConfigPath() string {
	return filepath.Join(GetAppDataDir(), "boccho.config.json")
}

func GetDefaultConfig() Config {
	return Config{
		FramesPath: getDefaultFramesPath(),
	}
}

func LoadConfig() (Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := GetDefaultConfig()
			if saveErr := SaveConfig(cfg); saveErr != nil {
				fmt.Printf("Warning: Could not save default config: %v\n", saveErr)
			}
			return cfg, nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.FramesPath == "" {
		cfg.FramesPath = getDefaultFramesPath()
	}

	return cfg, nil
}

func SaveConfig(cfg Config) error {
	configPath := GetConfigPath()

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func EnsureFramesDir(cfg Config) error {
	if _, err := os.Stat(cfg.FramesPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.FramesPath, 0755); err != nil {
			return fmt.Errorf("failed to create Frames directory: %w", err)
		}
		fmt.Printf("Created Frames folder: %s\n", cfg.FramesPath)
	}
	return nil
}
