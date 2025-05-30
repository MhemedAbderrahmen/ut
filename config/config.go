package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2" // Or encoding/json
)

type Config struct {
	AppName   string `yaml:"appname"`
	SecretKey string `yaml:"secretkey"`
}

var (
	cachedConfig *Config
	configMutex  sync.Mutex
)

func LoadConfig() (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if cachedConfig != nil {
		return cachedConfig, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("unable to find home directory: %w", err)
	}

	configFile := filepath.Join(home, ".ut-cli", "config.yml") // Or config.json

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("configuration file not found. Run 'ut config' to set it up")
		}
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg) // Or json.Unmarshal
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	cachedConfig = &cfg
	return cachedConfig, nil
}
