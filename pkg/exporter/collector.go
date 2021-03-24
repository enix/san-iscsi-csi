package exporter

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	csiRPCCall         *prometheus.CounterVec
	csiRPCCallDuration *prometheus.CounterVec
}

const (
	csiRPCCallMetric = "dothill_csi_rpc_call"
	csiRPCCallHelp   = "How many CSI RPC calls have been executed"

	csiRPCCallDurationMetric = "dothill_csi_rpc_call_duration"
	csiRPCCallDurationHelp   = "The total duration of CSI RPC calls"
)

func NewCollector() *Collector {
	return &Collector{
		csiRPCCall: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: csiRPCCallMetric,
				Help: csiRPCCallHelp,
			},
			[]string{"endpoint", "success"},
		),
		csiRPCCallDuration: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: csiRPCCallDurationMetric,
				Help: csiRPCCallDurationHelp,
			},
			[]string{"endpoint"},
		),
	}
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	collector.csiRPCCall.Describe(ch)
	collector.csiRPCCallDuration.Describe(ch)
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	collector.csiRPCCall.Collect(ch)
	collector.csiRPCCallDuration.Collect(ch)
}

func (collector *Collector) IncCSIRPCCall(method string, success bool) {
	collector.csiRPCCall.WithLabelValues(method, fmt.Sprintf("%t", success)).Inc()
}

func (collector *Collector) AddCSIRPCCallDuration(method string, duration time.Duration) {
	collector.csiRPCCallDuration.WithLabelValues(method).Add(float64(duration.Nanoseconds()) / 1000 / 1000 / 1000)
}
