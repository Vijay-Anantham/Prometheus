package main

import (
	"fmt"
	"net/http"

	"dopemeth/services"

	"github.com/gorilla/mux"
)

func main() {
	services.Wg.Add(1)
	go services.FetchStockPricesPeriodically()

	router := mux.NewRouter()
	router.HandleFunc("/api", services.GetStockPrice).Methods("GET")

	fmt.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
	// wait till all go routine gets finished
	services.Wg.Wait()
}
