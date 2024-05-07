package exporter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type ExporterOptions struct {
	Target                string
	IncludeDotDirectories bool
	Interval              time.Duration
}

func StartExporter(ctx context.Context, opts ExporterOptions, gauge *prometheus.GaugeVec) {
	for {
		select {
		case <-time.After(opts.Interval):
			directories, err := listDirectories(opts)
			if err != nil {
				fmt.Printf("Could not list directories: %v", err)
				break
			}
			for _, directory := range directories {
				size, err := dirSize(filepath.Join(opts.Target, directory))
				if err != nil {
					fmt.Printf("Could not calculate size of directory %s: %v", directory, err)
					break
				}
				gauge.With(prometheus.Labels{"persistentvolume": directory}).Set(size)
			}
		case <-ctx.Done():
			return
		}
	}
}

func listDirectories(opts ExporterOptions) ([]string, error) {
	directories := []string{}
	entries, err := os.ReadDir(opts.Target)
	if err != nil {
		fmt.Printf("Could read directories: %v", err)
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if !opts.IncludeDotDirectories && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}

// Func to calculate the size of a directory
func dirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return float64(size), nil
}
