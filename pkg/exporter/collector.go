package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	csiRPCCallCounters map[string]prometheus.Counter
}

var (
	csiRPCCallMetric = "dothill_csi_rpc_call"
	csiRPCCallHelp   = "How many CSI RPC calls have been executed"
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
