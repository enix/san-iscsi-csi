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
