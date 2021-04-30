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
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
)

// Exporter : Configuration (from command-line)
type Exporter struct {
	Port      int
	Collector *Collector

	listener   net.Listener
	server     *http.Server
	collectors []prometheus.Collector
}

func New(port int) *Exporter {
	exporter := &Exporter{
		Port:      port,
		Collector: NewCollector(),
	}
	exporter.RegisterCollector(exporter.Collector)
	return exporter
}

// ListenAndServe : Convenience function to start exporter
func (exporter *Exporter) ListenAndServe() error {
	if err := exporter.Listen(); err != nil {
		return err
	}

	return exporter.Serve()
}

// Listen : Listen for requests
func (exporter *Exporter) Listen() error {
	for _, collector := range exporter.collectors {
		err := prometheus.Register(collector)
		if err != nil {
			if registered, ok := err.(prometheus.AlreadyRegisteredError); ok {
				prometheus.Unregister(registered.ExistingCollector)
				prometheus.MustRegister(collector)
			}
		}
	}

	listen := fmt.Sprintf(":%d", exporter.Port)
	klog.Infof("listening on %s", listen)

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	exporter.listener = listener
	return nil
}

// Serve : Actually reply to requests
func (exporter *Exporter) Serve() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	exporter.server = &http.Server{
		Handler: mux,
	}

	return exporter.server.Serve(exporter.listener)
}

// Shutdown : Properly tear down server
func (exporter *Exporter) Shutdown() error {
	return exporter.server.Shutdown(context.Background())
}

func (exporter *Exporter) RegisterCollector(collector prometheus.Collector) {
	exporter.collectors = append(exporter.collectors, collector)
}
