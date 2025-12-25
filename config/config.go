package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var (
	ConfigPath = getConfigPath()
)

type Config struct {
	ProjectFolderPath string `toml:"path"`
}

// Loads and returns the config from the config file
func GetConfig() Config {
	initConfig()
	var config Config

	if _, err := toml.DecodeFile(ConfigPath, &config); err != nil {
		fmt.Println("Error reading config file:", err)
		return Config{}
	}
	return config
}

// Saves the config to the config file
func SaveConfig(config Config) error {
	initConfig()
	conf, err := os.OpenFile(ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer conf.Close()
	return toml.NewEncoder(conf).Encode(config)
}

// Makes sure a config file is present
func initConfig() error {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		file, err := os.Create(ConfigPath)
		if err != nil {
			return err
		}
		defer file.Close()
		SaveConfig(Config{""})
	}
	return nil
}

// Returns the correct path to the `config.toml` file by checking environment variables
func getConfigPath() string {
	switch true {
	case os.Getenv("DEVELOPMENT") == "true" || os.Getenv("DEV") == "true":
		return "./config.toml"
	default:
		path, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Error getting user config directory:", err)
			os.Exit(1)
		}
		return filepath.Join(path, "project-cli/project", "config.toml")
	}
}
