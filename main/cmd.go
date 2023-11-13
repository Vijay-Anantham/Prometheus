package main

import (
	"fmt"
	"net/http"

	"github.com/dopemeth/services"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api", services.GetStockPrice).Methods("GET")

	fmt.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
}
