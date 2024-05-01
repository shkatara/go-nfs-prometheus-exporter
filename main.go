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
	flag.StringVar(&target, "target-dir", "./", "Directory to scrape metrics from")
	fmt.Println("Reading from", target, "and exposing metrics at 127.0.0.1:8000/metrics")
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8000", nil)
	dir := utils.FindDir(target)
	for {
		time.Sleep(5 * time.Second)
		for _, data := range dir {
			size := utils.DirSize(data)
			gauge.With(prometheus.Labels{"persistentvolume": data}).Set(float64(size))
		}
	}
}
