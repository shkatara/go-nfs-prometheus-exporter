package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	flag.StringVar(&target, "target-dir", "", "Directory to scrape metrics from")
	flag.StringVar(&listenAddress, "listen-address", "0.0.0.0", "Listen address")
	flag.IntVar(&listenPort, "listen-port", 8000, "Listen port")
	flag.BoolVar(&includeDotDirectories, "include-dot-dirs", false, "Include directories starting with a dot")

	flag.Parse()

	opts := exporter.ExporterOptions{
		Target:                target,
		IncludeDotDirectories: includeDotDirectories,
	}

	listenAddress := fmt.Sprintf("%s:%v", listenAddress, listenPort)
	fmt.Printf("Reading from %s and serving metrics at %s/metrics", target, listenAddress)

	http.Handle("/metrics", promhttp.Handler())

	// Remove this because if a new persistent volume is added, it will not be scraped
	// since we discover only at the start of the program with this approach
	// dir := utils.FindDir(target)

	err := http.ListenAndServe(listenAddress, nil)
	if err != nil {
		return err
	}

	// We start this asynchronously so that we can serve metrics in the background
	go exporter.StartExporter(ctx, opts, gauge)

	// Block and wait for interrupt.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	return nil
}
