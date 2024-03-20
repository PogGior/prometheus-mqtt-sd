package util

import (
	"gopkg.in/yaml.v2"
	"os"
)

type SdConfig struct {
	Address  string `yaml:"address"`
	Topic    string `yaml:"topic"`
	ClientID string `yaml:"client_id"`
}

func LoadConfig(file string) (*SdConfig, error) {
	var config SdConfig

	configFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
