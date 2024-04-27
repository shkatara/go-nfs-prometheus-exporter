package main

import (
	"fmt"
	"net/http"
	"time"

	"example.com/nfs-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "dir_size_bytes",
		Help:        "The size of the directory in bytes",
		ConstLabels: prometheus.Labels{"job": "go-nfs-exporter"},
	})
)

func main() {
	var target string = "./"
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8000", nil)
	dir := utils.FindDir(target)
	for {
		time.Sleep(5 * time.Second)
		for _, data := range dir {
			size := utils.DirSize(data)
			gauge.Set(float64(size))
			fmt.Println(size)
		}
	}
}
