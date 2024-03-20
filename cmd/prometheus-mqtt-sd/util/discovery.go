package util

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)


type discovery struct {
	topic  string
	client mqtt.Client
}


func (config *SdConfig) InitDiscovery() (*discovery, error) {
	client, err := config.initClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cd := &discovery{
		topic:  config.Topic,
		client: client,
	}
	return cd, nil
}

func (d *discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	// Create a channel to handle incoming MQTT messages
	msgChan := make(chan *mqtt.Message)

	// Set up the MQTT message handler
	d.client.Subscribe(d.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		msgChan <- &msg
	})

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			// Process the MQTT message to create a targetgroup.Group
			tg, err := d.processMessage(msg)
			if err != nil {
				continue
			}

			// Send the target group to the provided channel
			ch <- tg
		}
	}
}

func (d *discovery) processMessage(msg *mqtt.Message) ([]*targetgroup.Group, error) {
	// TODO: Implement this function to process the MQTT message and create a targetgroup.Group
	// You will likely need to parse the message payload and possibly use additional metadata
	// from the MQTT message.

	g := []struct {
		Targets []string          `json:"targets"`
		Labels  map[string]string `json:"labels"`
	}{}

	// deserialize the message payload into the struct g
	if err := json.Unmarshal((*msg).Payload(), &g); err != nil {
		return nil, errors.WithStack(err)
	}

	println("Processing message: ", g)
	println("message lenght: ", len(g) )
	tgs := make([]*targetgroup.Group, 0, len(g))

	for _, group := range g {
		targets := make([]model.LabelSet, 0, len(group.Targets))
		for _, t := range group.Targets {
			targets = append(targets, model.LabelSet{
				model.AddressLabel: model.LabelValue(t),
			})
		}
		labels := make(model.LabelSet, len(group.Labels))
		for k, v := range group.Labels {
			labels[model.LabelName(k)] = model.LabelValue(v)
		}
		if err := labels.Validate(); err != nil {
			return nil, errors.WithStack(err)
		}

		tgs = append(tgs, &targetgroup.Group{
			Targets: targets,
			Labels:  labels,
			Source: group.Labels["job"],
		})
		println("tgs: ", tgs)
		println("tgs lenght: ", len(tgs) )
		println("targets: ", tgs[0].Targets)
	}
	return tgs, nil

}