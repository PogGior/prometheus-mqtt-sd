package util

import (
    "os"

    "gopkg.in/yaml.v2"
)

type SdConfig struct {
    Address string `yaml:"address"`
    Topic   string `yaml:"topic"`
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