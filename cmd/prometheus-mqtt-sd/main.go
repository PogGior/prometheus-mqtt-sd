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

	cfg, err := util.LoadConfig(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(1)	
	}

	disc, err := cfg.InitDiscovery(logger)
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
