package main

import (
	"fmt"
	"net/http"

	services "github.com/dopemeth/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ping(w http.ResponseWriter, req *http.Request) {
	pingCounter.Inc()
	fmt.Fprintf(w, "Request count incremented")
}

var pingCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ping_request_count",
		Help: "No of request handled by Ping handler",
	},
)

func main() {
	prometheus.MustRegister(pingCounter)
	http.HandleFunc("/api", services.ApiHandler) // Define the API endpoint and its corresponding handler function
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8090", nil)
}
