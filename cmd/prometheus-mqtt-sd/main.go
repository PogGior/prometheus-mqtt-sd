// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	prom_discovery "github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/documentation/examples/custom-sd/adapter"
)

var (
	a             = kingpin.New("sd adapter usage", "Tool to generate file_sd target files for unimplemented SD mechanisms.")
	outputFile    = a.Flag("output.file", "Output file for file_sd compatible file.").Default("custom_sd.json").String()
	listenAddress = a.Flag("listen.address", "The address of mqtt broker.").Default("localhost:8500").String()
	topic         = a.Flag("topic", "The topic to subscribe to.").Default("example").String()
	logger        log.Logger
)

// Note: create a config struct for your custom SD type here.
type sdConfig struct {
	address string
	topic   string
	logger  log.Logger
}

type discovery struct {
	// ...
	topic  string
	client mqtt.Client
	logger log.Logger
}

// ...

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
				level.Error(d.logger).Log("msg", "Error processing MQTT message", "err", err)
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
		})
	}
	return tgs, nil

}

func newDiscovery(conf sdConfig) (*discovery, error) {
	// create a new MQTT client
	opts := mqtt.NewClientOptions().AddBroker(conf.Address)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	cd := &discovery{
		topic:  *topic,
		client: client,
		logger: logger,
	}
	return cd, nil
}

func main() {
	a.HelpFlag.Short('h')

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	logger = log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	ctx := context.Background()

	// NOTE: create an instance of your new SD implementation here.
	cfg := sdConfig{
		address: *listenAddress,
		topic:   *topic,
		logger:  logger,
	}

	disc, err := newDiscovery(cfg)
	if err != nil {
		fmt.Println("err: ", err)
	}

	if err != nil {
		level.Error(logger).Log("msg", "failed to create discovery metrics", "err", err)
		os.Exit(1)
	}

	reg := prometheus.NewRegistry()
	refreshMetrics := prom_discovery.NewRefreshMetrics(reg)
	metrics, err := prom_discovery.RegisterSDMetrics(reg, refreshMetrics)
	if err != nil {
		level.Error(logger).Log("msg", "failed to register service discovery metrics", "err", err)
		os.Exit(1)
	}

	sdAdapter := adapter.NewAdapter(ctx, *outputFile, "exampleSD", disc, logger, metrics, reg)
	sdAdapter.Run()

	<-ctx.Done()
}
