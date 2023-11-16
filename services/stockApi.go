package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO: Webex alert for BUY SELL

// Adding prometheus data
var (
	Pricecounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stock_price_requests_total",
			Help: "Total number of stock price requests.",
		},
		[]string{"symbol"},
	)
	Gaincounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stock_gain_incured",
			Help: "Number of time 5% Gain incured.",
		},
		[]string{"symbol"},
	)
	Losscounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stock_loss_incured",
			Help: "Number of time 5% loss incured.",
		},
		[]string{"symbol"},
	)

	PriceGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stock_price",
			Help: "Current stock price.",
		},
		[]string{"symbol"},
	)

	resultCh string
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

type Output struct {
	Name  string
	Price float64
}

func GetStockPrice(w http.ResponseWriter, r *http.Request) {
	_, err := FetchAndPrintStockPrices()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := resultCh
	w.Write([]byte(result))
}

// func FetchStockPricesPeriodically(w http.ResponseWriter, r *http.Request) {
// 	// defer Wg.Done()

// 	ticker := time.NewTicker(time.Second * 10)
// 	defer ticker.Stop()
// 	defer log.Print("got hit")

// 	for {
// 		select {
// 		case <-ticker.C:
// 			if err := fetchAndPrintStockPrices(w); err != nil {
// 				fmt.Println("Error fetching stock prices:", err)
// 			}
// 			fmt.Print("still in loop")
// 		}
// 	}
// }

func FetchAndPrintStockPrices() (Output, error) {
	log.Print("Entered function")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=1min&apikey=%s", stockSymbol, apiKey)
	Pricecounter.WithLabelValues(stockSymbol).Inc()
	log.Print("Updated counter")
	resp, err := http.Get(url)
	if err != nil {
		return Output{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Output{}, err
	}

	var data IntradayResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return Output{}, err
	}

	latestData := getLatestPrice(data.TimeSeries)
	result, _ := strconv.ParseFloat(latestData.Close, 64)
	fmt.Print("Fetching stock prices: \n")
	response := fmt.Sprintf(`{"Stock name": "%s", "latest_price": "%s"}`, stockSymbol, latestData.Close)
	fmt.Println("Response:", string(response))
	pollerout := Output{
		Name:  stockSymbol,
		Price: result,
	}
	log.Print("Output ready")
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(response))
	PriceGauge.WithLabelValues(stockSymbol).Set(result)
	log.Print("Gauge value set")
	resultCh = response
	log.Print("Funtion terminating gracefully")
	return pollerout, nil
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
