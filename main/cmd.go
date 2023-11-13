package main

import (
	"fmt"
	"net/http"

	"dopemeth/services"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(services.Pricecounter)
	prometheus.MustRegister(services.Gaincounter)
	prometheus.MustRegister(services.Losscounter)
	prometheus.MustRegister(services.PriceGauge)

}

func main() {
	services.Wg.Add(1)
	go services.FetchStockPricesPeriodically()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/api", services.GetStockPrice).Methods("GET")

	fmt.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
	// wait till all go routine gets finished
	services.Wg.Wait()
}
