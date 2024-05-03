package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	flag.StringVar(&target, "target-dir", target, "Directory to scrape metrics from")
	flag.Parse()
	outboundIP := GetOutboundIP()
	listenAddress := fmt.Sprintf("%s:8000", outboundIP)
	fmt.Print(fmt.Sprintf("Reading from %s and serving metrics at %s:8000/metrics", target, outboundIP))
	http.Handle("/metrics", promhttp.Handler())
	dir := utils.FindDir(target)
	go http.ListenAndServe(listenAddress, nil)
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
