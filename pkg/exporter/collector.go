/*
 * Copyright (c) 2021 Enix, SAS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * Authors:
 * Paul Laffitte <paul.laffitte@enix.fr>
 * Arthur Chaloin <arthur.chaloin@enix.fr>
 * Alexandre Buisine <alexandre.buisine@enix.fr>
 */

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
