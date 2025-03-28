package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

/* ================================== TYPES ================================= */

type TimerConfig struct {
	DurationMinutes int `yaml:"duration_minutes"`
	BreakSeconds    int `yaml:"break_seconds"`
}

type NotificationConfig struct {
	UseAlert bool `yaml:"use_alert"`
}

type AppConfig struct {
	Timer         TimerConfig        `yaml:"timer"`
	Notifications NotificationConfig `yaml:"notifications"`
}

/* ================================ VARIABLES =============================== */

const APP_NAME string = "look_away"
const CONFIG_FILE_NAME string = "config.yaml"

/* ================================ FUNCTIONS =============================== */

func LoadConfig() (*AppConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("could not get config path: %v", err)
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Config file not found. Creating default config...")
		err := createDefaultConfig(path)
		if err != nil {
			return nil, fmt.Errorf("error creating default config: %v", err)
		}
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %v", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("failed to parse yaml %v", err)
	}

	return &config, nil
}

func GetConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config directory not found: %v", err)
	}

	configPath := filepath.Join(userConfigDir, APP_NAME, CONFIG_FILE_NAME)
	return configPath, nil
}

func createDefaultConfig(configPath string) error {
	defaultConfig := AppConfig{
		Timer: TimerConfig{
			DurationMinutes: 20,
			BreakSeconds:    20,
		},
		Notifications: NotificationConfig{
			UseAlert: true,
		},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) GetTimerDuration() time.Duration {
	return time.Duration(c.Timer.DurationMinutes) * time.Minute
}

func (c *AppConfig) GetBreakSeconds() time.Duration {
	return time.Duration(c.Timer.BreakSeconds) * time.Second
}
