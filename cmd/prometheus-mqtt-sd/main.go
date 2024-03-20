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
	"fmt"
	"os"
	"prometheus-mqtt-sd/cmd/prometheus-mqtt-sd/util"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	prom_discovery "github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/documentation/examples/custom-sd/adapter"
)

var (
	a             = kingpin.New("sd adapter usage", "Tool to generate file_sd target files for unimplemented SD mechanisms.")
	outputFile    = a.Flag("output.file", "Output file for file_sd compatible file.").Default("custom_sd.json").String()
	configFile    = a.Flag("config.file", "The path of the configuration file").Default("config.yaml").String()
	logger        log.Logger
)

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

	// Load the configuration from the file
	cfg, err := util.LoadConfig(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)	
	}

	disc, err := cfg.InitDiscovery()
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
