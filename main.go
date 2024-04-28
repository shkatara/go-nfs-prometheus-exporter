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
	gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dir_size_bytes",
			Help: "Directory size in bytes.",
		},
		[]string{"persistentvolumes"})
)

var target string = "./"

func main() {
	fmt.Println("Starting server...")
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8000", nil)
	dir := utils.FindDir(target)
	for {
		time.Sleep(5 * time.Second)
		for _, data := range dir {
			size := utils.DirSize(data)
			gauge.With(prometheus.Labels{"persistentvolumes": data}).Set(float64(size))
		}
	}
}
