package util

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


type discovery struct {
	// ...
	topic  string
	client mqtt.Client
	logger log.Logger
}


func (config *SdConfig)InitDiscovery() (*discovery, error) {
	opts := mqtt.NewClientOptions().AddBroker(config.Address)
	opts.SetClientID("prometheus-mqtt-sd")
	// Create a new MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	cd := &discovery{
		topic:  config.Topic,
		client: client,
	}
	return cd, nil
}