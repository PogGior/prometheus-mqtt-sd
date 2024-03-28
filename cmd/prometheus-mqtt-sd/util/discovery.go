package util

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type discovery struct {
	topic  string
	client mqtt.Client
	logger log.Logger
}

func (discovery *discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {

	brokerMessageChain := discovery.brokerSubscribe()

	for {
		select {

		case <-ctx.Done():
			return

		case msg := <-brokerMessageChain:
			level.Debug(discovery.logger).Log("msg", "received message", "payload", string((*msg).Payload()))
			tg, err := ProcessMessage((*msg).Payload())
			if err != nil {
				continue
			}
			ch <- tg

		}
	}
}

func (discovery *discovery) brokerSubscribe() chan *mqtt.Message {

	brokerMessageChain := make(chan *mqtt.Message)

	discovery.client.Subscribe(discovery.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		brokerMessageChain <- &msg
	})

	return brokerMessageChain
}

func ProcessMessage(payload []byte) ([]*targetgroup.Group, error) {

	group := []struct {
		Targets []string          `json:"targets"`
		Labels  map[string]string `json:"labels"`
		Source  string            `json:"source"`
	}{}

	if err := json.Unmarshal(payload, &group); err != nil {
		return nil, errors.WithStack(err)
	}

	targetGroups := make([]*targetgroup.Group, 0, len(group))

	for _, group := range group {

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

		targetGroups = append(targetGroups, &targetgroup.Group{
			Targets: targets,
			Labels:  labels,
			Source:  group.Source,
		})

	}

	return targetGroups, nil
}
