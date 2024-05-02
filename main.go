package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"example.com/nfs-exporter/utils"
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
	target string
)

func main() {
	flag.StringVar(&target, "target-dir", target, "Directory to scrape metrics from")
	flag.Parse()
	fmt.Println("Reading from", target, "and exposing metrics at 127.0.0.1:8000/metrics")

	http.Handle("/metrics", promhttp.Handler())
	dir := utils.FindDir(target)
	go http.ListenAndServe(":8000", nil)
	for {
		time.Sleep(30 * time.Second)
		for _, data := range dir {
			fmt.Println("Scraping data from", data)
			size := utils.DirSize(target, data)
			fmt.Println(data, ":", size, "bytes")
			gauge.With(prometheus.Labels{"persistentvolume": data}).Set(size)
		}
	}
}
