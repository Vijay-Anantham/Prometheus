package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiKey      = "TOB0OSLJF393DSV0"
	stockSymbol = "CSCO" // Replace with the stock symbol you want to monitor
	interval    = 1      // Interval in minutes
)

type IntradayResponse struct {
	MetaData   map[string]string `json:"Meta Data"`
	TimeSeries map[string]Price  `json:"Time Series (1min)"`
}

type Price struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

func GetStockPrice(w http.ResponseWriter, r *http.Request) {
	err := fetchAndPrintStockPrices(w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func fetchAndPrintStockPrices(w http.ResponseWriter) error {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=1min&apikey=%s", stockSymbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data IntradayResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	latestData := getLatestPrice(data.TimeSeries)
	fmt.Sprint("Fetching stock prices: \n")
	response := fmt.Sprintf(`{"Stock name": "%s", "latest_price": "%s"}`, stockSymbol, latestData.Close)
	fmt.Println("Response:", string(response))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))

	return nil
}

func getLatestPrice(timeSeries map[string]Price) Price {
	// Extract the latest stock price from the response
	var latestData Price
	for _, v := range timeSeries {
		latestData = v
		break
	}
	return latestData
}
