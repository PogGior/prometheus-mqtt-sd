package util

import (
	"os"

	"github.com/go-kit/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
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

func (config *SdConfig) InitDiscovery(logger log.Logger) (*discovery, error) {

	client, err := config.initClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cd := &discovery{
		topic:  config.Topic,
		client: client,
		logger: logger,
	}

	return cd, nil
}

func (config *SdConfig)initClient() (mqtt.Client,error) {

	opts := mqtt.NewClientOptions().AddBroker(config.Address)
	opts.SetClientID(config.ClientID)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client,nil
}
