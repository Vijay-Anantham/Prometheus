package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type StockPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	// Make an API call to Alpha Vantage
	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=IBM&apikey=TOB0OSLJF393DSV0")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the API response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON response
	var data struct {
		GlobalQuote struct {
			Symbol string `json:"01. symbol"`
			Price  string `json:"05. price"`
		} `json:"Global Quote"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Create a StockPrice struct with the parsed data
	stockPrice := StockPrice{
		Symbol: data.GlobalQuote.Symbol,
	}
	fmt.Sscanf(data.GlobalQuote.Price, "%f", &stockPrice.Price)

	// Convert the StockPrice struct to JSON
	jsonData, err := json.Marshal(stockPrice)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Content-Type header and send the JSON response back to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
