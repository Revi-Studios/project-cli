package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var ConfigPath = getConfigPath()

type Config struct {
	Debug    bool     `toml:"debug"`
	Editor   string   `toml:"editor"`
	Projects Projects `toml:"projects"`
	Defaults Defaults `toml:"defaults"`
}

type Projects struct {
	Path      string `toml:"path"`
	DebugPath string `toml:"debug_path"`
}

type Defaults struct {
	Tag  string `toml:"tag"`
	View string `toml:"view"`
	Open string `toml:"open"`
}

func (this *Config) ProjectPath() string {
	if this.Debug {
		return this.Projects.DebugPath
	}
	return this.Projects.Path
}

// Loads and returns the config from the config file
func GetConfig() (*Config, error) {
	if err := validateConfigFile(); err != nil {
		return nil, fmt.Errorf("validating the config: %w", err)
	}
	var config *Config

	if _, err := toml.DecodeFile(ConfigPath, &config); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	return config, nil
}

// Saves the config to the config file
func SaveConfig(config *Config) error {
	validateConfigFile()
	conf, err := os.OpenFile(ConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer conf.Close()
	return toml.NewEncoder(conf).Encode(*config)
}

// Makes sure a config file is present
func validateConfigFile() error {
	if err := os.MkdirAll(filepath.Dir(ConfigPath), 0o755); err != nil {
		return fmt.Errorf("creating directories %q: %w", ConfigPath, err)
	}
	conf, err := os.OpenFile(ConfigPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}

		return fmt.Errorf("opening the file %v: %w", ConfigPath, err)
	}
	defer conf.Close()
	return nil
}

// Returns the correct path to the `config.toml` file by checking environment variables
func getConfigPath() string {
	switch true {

	case os.Getenv("DEVELOPMENT") == "true" || os.Getenv("DEV") == "true":
		return "./config.toml"
	case os.Getenv("CONFIG_PATH") != "":
		return os.Getenv("CONFIG_PATH")
	default:
		path, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Error getting user config directory:", err)
			os.Exit(1)
		}
		return filepath.Join(path, "project", "config.toml")
	}
}
