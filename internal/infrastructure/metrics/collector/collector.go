package collector

import (
	"math/rand/v2"
	"runtime"
)

type RuntimeMetrics struct {
	memStats  runtime.MemStats
	pollCount int64
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{}
}

func (rm *RuntimeMetrics) Collect() map[string]interface{} {
	runtime.ReadMemStats(&rm.memStats)
	rm.pollCount++

	return map[string]interface{}{
		// Gauge metrics
		"Alloc":         float64(rm.memStats.Alloc),
		"BuckHashSys":   float64(rm.memStats.BuckHashSys),
		"Frees":         float64(rm.memStats.Frees),
		"GCCPUFraction": rm.memStats.GCCPUFraction,
		"GCSys":         float64(rm.memStats.GCSys),
		"HeapAlloc":     float64(rm.memStats.HeapAlloc),
		"HeapIdle":      float64(rm.memStats.HeapIdle),
		"HeapInuse":     float64(rm.memStats.HeapInuse),
		"HeapObjects":   float64(rm.memStats.HeapObjects),
		"HeapReleased":  float64(rm.memStats.HeapReleased),
		"HeapSys":       float64(rm.memStats.HeapSys),
		"LastGC":        float64(rm.memStats.LastGC),
		"Lookups":       float64(rm.memStats.Lookups),
		"MCacheInuse":   float64(rm.memStats.MCacheInuse),
		"MCacheSys":     float64(rm.memStats.MCacheSys),
		"MSpanInuse":    float64(rm.memStats.MSpanInuse),
		"MSpanSys":      float64(rm.memStats.MSpanSys),
		"Mallocs":       float64(rm.memStats.Mallocs),
		"NextGC":        float64(rm.memStats.NextGC),
		"NumForcedGC":   float64(rm.memStats.NumForcedGC),
		"NumGC":         float64(rm.memStats.NumGC),
		"OtherSys":      float64(rm.memStats.OtherSys),
		"PauseTotalNs":  float64(rm.memStats.PauseTotalNs),
		"StackInuse":    float64(rm.memStats.StackInuse),
		"StackSys":      float64(rm.memStats.StackSys),
		"Sys":           float64(rm.memStats.Sys),
		"TotalAlloc":    float64(rm.memStats.TotalAlloc),
		"RandomValue":   rand.Float64(),

		// Counter metrics
		"PollCount": rm.pollCount,
	}
}
