package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

var (
	ErrConfigNotFound = errors.New("configuration file not found")
	ErrAPIKeyMissing  = errors.New("api key not set")
)

func LoadConfig() (*Config, error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	if cachedConfig != nil {
		return cachedConfig, nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("unable to find home directory: %w", ErrConfigNotFound)
	}

	configFile := filepath.Join(home, ".ut-cli", "config.yml") // Or config.json

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg) // Or json.Unmarshal
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}
	if strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, ErrAPIKeyMissing
	}

	cachedConfig = &cfg
	return cachedConfig, nil
}
