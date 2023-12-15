package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port     string   `yaml:"Port"`
	Services []string `yaml:"Services"`
}

func LoadConfig() (*Config, error) {
	path := "./config/config.yaml"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	var cfg Config
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
