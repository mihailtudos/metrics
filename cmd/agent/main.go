package main

import (
	"time"

	"github.com/mihailtudos/metrics/internal/infrastructure/metrics/collector"
	"github.com/mihailtudos/metrics/internal/infrastructure/metrics/reporter"
)

const ServerURL = "http://localhost:8080"

var (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

func main() {
	reporter := reporter.NewMetricsReporter(ServerURL)
	metrics := collector.NewRuntimeMetrics()
	pollTicker := time.NewTicker(pollInterval)
	reportTicker := time.NewTicker(reportInterval)

	for {
		select {
		case <-pollTicker.C:
			_ = metrics.Collect()
		case <-reportTicker.C:
			reporter.ReportMetrics(metrics.Collect())
		}
	}
}
