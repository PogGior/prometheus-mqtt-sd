package util

import mqtt "github.com/eclipse/paho.mqtt.golang"


func (config *SdConfig)initClient() (mqtt.Client,error) {
	opts := mqtt.NewClientOptions().AddBroker(config.Address)
	opts.SetClientID(config.ClientID)
	// Create a new MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client,nil
}