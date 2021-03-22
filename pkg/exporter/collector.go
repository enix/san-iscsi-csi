package exporter

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	csiRPCCallCounters map[string]prometheus.Counter
}

const (
	csiRPCCallMetric = "dothill_csi_rpc_call"
	csiRPCCallHelp   = "How many CSI RPC calls have been executed"

	csiRPCCallDurationMetric = "dothill_csi_rpc_call_duration"
	csiRPCCallDurationHelp   = "The total duration of CSI RPC calls"
)

func NewCollector() *Collector {
	return &Collector{
		csiRPCCallCounters: map[string]prometheus.Counter{},
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, counter := range collector.csiRPCCallCounters {
		ch <- counter.Desc()
	}
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, counter := range collector.csiRPCCallCounters {
		ch <- counter
	}
}

func (collector *Collector) getCSIRPCCallCounter(method string, success bool) prometheus.Counter {
	ID := method + ":" + fmt.Sprintf("%t", success)

	if counter, ok := collector.csiRPCCallCounters[ID]; ok {
		return counter
	}

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: csiRPCCallMetric,
		Help: csiRPCCallHelp,
		ConstLabels: prometheus.Labels{
			"method":  method,
			"success": fmt.Sprintf("%t", success),
		},
	})
	collector.csiRPCCallCounters[ID] = counter

	return counter
}

func (collector *Collector) IncCSIRPCCall(method string, success bool) {
	collector.getCSIRPCCallCounter(method, success).Inc()
}

func (collector *Collector) getCSIRPCCallDurationCounter(method string) prometheus.Counter {
	if counter, ok := collector.csiRPCCallCounters[method]; ok {
		return counter
	}

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: csiRPCCallDurationMetric,
		Help: csiRPCCallDurationHelp,
		ConstLabels: prometheus.Labels{
			"method": method,
		},
	})
	collector.csiRPCCallCounters[method] = counter

	return counter
}

func (collector *Collector) AddCSIRPCCallDuration(method string, duration time.Duration) {
	collector.getCSIRPCCallDurationCounter(method).Add(float64(duration.Nanoseconds()) / 1000 / 1000 / 1000)
}
