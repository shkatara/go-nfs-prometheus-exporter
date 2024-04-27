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

func recordMetrics(size int64) {
	byteSize.Set(float64(size))
}

var (
	byteSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dir_size_bytes",
		Help: "The size of the directory in bytes",
	})
)

func main() {
	var target string = "./"
	fmt.Println("Server started on port 8000")
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8000", nil)

	dir := utils.FindDir(target)
	for {
		time.Sleep(5 * time.Second)
		for _, data := range dir {
			size := utils.DirSize(data)
			recordMetrics(size)
			fmt.Println(size)
		}
	}
}
