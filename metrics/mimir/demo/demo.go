package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "duration",
			Subsystem: "demo",
			Help:      "API latency distributions.",
		},
		[]string{"app"},
	)
)

var promHandler = promhttp.HandlerFor(
	prometheus.DefaultGatherer,
	promhttp.HandlerOpts{},
)

func init() {
	prometheus.MustRegister(apiDurationHistogram)
}

func main() {
	go func() {
		err := http.ListenAndServe(":8080", promHandler)
		fmt.Println(err)
	}()
	rand.Seed(time.Now().UnixNano())
	for {
		apiDurationHistogram.With(map[string]string{"app": "demo"}).Observe(rand.Float64())
		time.Sleep(2 * time.Second)
	}
}
