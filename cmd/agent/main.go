package main

import (
	"time"

	"github.com/mihailtudos/metrics/internal/infrastructure/config/agent"
	"github.com/mihailtudos/metrics/internal/infrastructure/metrics/collector"
	"github.com/mihailtudos/metrics/internal/infrastructure/metrics/reporter"
)

func main() {
	cfgAgent := agent.NewAgentConfig()

	reporter := reporter.NewMetricsReporter(cfgAgent.ServerAddress)
	metrics := collector.NewRuntimeMetrics()
	pollTicker := time.NewTicker(cfgAgent.PollInterval)
	reportTicker := time.NewTicker(cfgAgent.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			_ = metrics.Collect()
		case <-reportTicker.C:
			reporter.ReportMetrics(metrics.Collect())
		}
	}
}
