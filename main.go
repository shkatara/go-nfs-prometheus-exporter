package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"nfs-exporter/exporter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "persistentvolume_used_bytes",
			Help: "Persistent Volume Used in bytes.",
		},
		[]string{"persistentvolume"})
)

var (
	target                string
	listenAddress         string
	listenPort            int
	includeDotDirectories bool
	interval              string
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	flag.StringVar(&target, "target-dir", "/Users/jmotz/Code/Github/", "Directory to scrape metrics from")
	flag.StringVar(&listenAddress, "listen-address", "0.0.0.0", "Listen address")
	flag.IntVar(&listenPort, "listen-port", 8000, "Listen port")
	flag.BoolVar(&includeDotDirectories, "include-dot-dirs", false, "Include directories starting with a dot")
	flag.StringVar(&interval, "interval", "30s", "Interval to scrape directories (e.g. 30s, 1m, 1h)")

	flag.Parse()

	parsedInterval, err := time.ParseDuration(interval)
	if err != nil {
		fmt.Printf("Could not parse interval: %v", err)
		return err
	}

	opts := exporter.ExporterOptions{
		Target:                target,
		IncludeDotDirectories: includeDotDirectories,
		Interval:              parsedInterval,
	}

	// We start this asynchronously so that we can serve metrics in the background
	go exporter.StartExporter(ctx, opts, gauge)

	// Start Prometheus metrics server
	listenAddress := fmt.Sprintf("%s:%v", listenAddress, listenPort)
	fmt.Printf("Reading from %s and serving metrics at %s/metrics", target, listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listenAddress, nil)
}
