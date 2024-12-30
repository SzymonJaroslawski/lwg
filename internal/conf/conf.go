package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SzymonJaroslawski/lwg/internal/utils"
	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName = "config.yaml"
	ConfigPerms    = 0755
)

type Config struct {
	Paths struct {
		MainDir string `yaml:"main_dir"`
		Games   string `yaml:"games"`
		Runners string `yaml:"runners"`
	} `yaml:"paths"`
	Preferences `yaml:"preferences"`
}

type Preferences struct {
	Runners struct {
		DefaultFriendlyName string    `yaml:"default_friendly_name"`
		DefaultPath         string    `yaml:"default_path"`
		DefaultID           uuid.UUID `yaml:"default_id"`
	} `yaml:"runners"`
}

// Returns new config with default configuration
func NewConfig() *Config {
	return &Config{
		Paths: struct {
			MainDir string "yaml:\"main_dir\""
			Games   string "yaml:\"games\""
			Runners string "yaml:\"runners\""
		}{
			MainDir: filepath.Join(xdg.ConfigHome, "lwg"),
			Games:   filepath.Join(xdg.ConfigHome, "lwg", "games"),
			Runners: filepath.Join(xdg.ConfigHome, "lwg", "runners"),
		},
		Preferences: Preferences{
			Runners: struct {
				DefaultFriendlyName string    "yaml:\"default_friendly_name\""
				DefaultPath         string    "yaml:\"default_path\""
				DefaultID           uuid.UUID "yaml:\"default_id\""
			}{
				DefaultFriendlyName: "",
				DefaultPath:         "",
				DefaultID:           uuid.Nil,
			},
		},
	}
}

// Argument `path` should be a path to directory containing config, not config.yaml
func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Path: %s, does not exist: %v", path, err)
	}

	if st, err := os.Stat(path); !st.IsDir() {
		if err != nil {
			return nil, fmt.Errorf("os.Stat unexpected error: %v", err)
		}

		return nil, fmt.Errorf("Path: %s, is not a directory", path)
	}

	confFilePath := filepath.Join(path, ConfigFileName)

	f, err := os.Open(confFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(conf)
	if err != nil {
		return nil, fmt.Errorf("Config file: %s, decode error: %v", confFilePath, err)
	}

	return conf, nil
}

func (c *Config) Save() (bool, error) {
	if !utils.Exsits(c.Paths.MainDir) {
		ok, err := c.init()
		if !ok {
			return false, fmt.Errorf("Failed to scaffold config directory: %s, %v", c.Paths.MainDir, err)
		}
	}

	configFilePath := filepath.Join(c.Paths.MainDir, ConfigFileName)

	if !utils.Exsits(configFilePath) {
		f, err := os.Create(configFilePath)
		if err != nil {
			return false, fmt.Errorf("Failed creating config file: %s, %v", configFilePath, err)
		}
		defer f.Close()

		encoder := yaml.NewEncoder(f)
		err = encoder.Encode(c)
		if err != nil {
			return false, fmt.Errorf("Failed to save config: %s, %v", configFilePath, err)
		}
	}

	ok, err := c.update()
	if !ok {
		return false, err
	}

	return true, nil
}

var ErrConfFileNotExsits = errors.New("Atempted to update but config file does not exist")

func (c *Config) update() (bool, error) {
	configFilePath := filepath.Join(c.Paths.MainDir, ConfigFileName)

	if !utils.Exsits(configFilePath) {
		return false, ErrConfFileNotExsits
	}

	err := os.Remove(configFilePath)
	if err != nil {
		return false, err
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(c)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Config) init() (bool, error) {
	err := os.MkdirAll(c.Paths.MainDir, ConfigPerms)
	if err != nil {
		return false, err
	}
	err = os.MkdirAll(c.Paths.Games, ConfigPerms)
	if err != nil {
		return false, err
	}
	err = os.MkdirAll(c.Paths.Runners, ConfigPerms)
	if err != nil {
		return false, err
	}

	return true, nil
}
