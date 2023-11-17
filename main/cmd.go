package main

import (
	"fmt"
	"net/http"

	"dopemeth/poller"
	"dopemeth/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: Try adding a input collection for name of stock to scrap

func init() {
	prometheus.MustRegister(services.Pricecounter)
	prometheus.MustRegister(services.Gaincounter)
	prometheus.MustRegister(services.Losscounter)
	prometheus.MustRegister(services.PriceGauge)

}

func main() {
	// services.Wg.Add(1)
	// go services.FetchStockPricesPeriodically()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/api", services.GetStockPrice).Methods("GET")
	// router.HandleFunc("/apiRealtime", services.FetchStockPricesPeriodically).Methods("GET")
	go func() {
		fmt.Println("Metrics server listening on :8080")
		http.ListenAndServe(":8080", router)
	}()

	// Start a Goroutine to periodically poll the API
	go poller.PollApi()

	// Run indefinitely to keep the program alive
	select {}
	// fmt.Println("API server listening on :8080")
	// http.ListenAndServe(":8080", router)
	// wait till all go routine gets finished
	// services.Wg.Wait()
}
