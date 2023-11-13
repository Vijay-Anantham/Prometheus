package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO: Add prometheus montitoring
// TODO: Add a counter 'Gaincounter' -> gain >= 5%
// TODO: Add a counter 'Losscounter' -> loss >= 5%
// TODO: Add a gauge 'fluctuations' --> lookup what gauge is
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

	resultCh = make(chan string, 1)
	Wg       sync.WaitGroup
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
	err := fetchAndPrintStockPrices()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := <-resultCh
	w.Write([]byte(result))
}

func FetchStockPricesPeriodically() {
	defer Wg.Done()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	defer log.Print("got hit")

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := fetchAndPrintStockPrices(); err != nil {
					fmt.Println("Error fetching stock prices:", err)
				}
				fmt.Print("still in loop")
			}
		}
	}()
	<-done
}

func fetchAndPrintStockPrices() error {
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
	result, _ := strconv.ParseFloat(latestData.Close, 64)
	fmt.Print("Fetching stock prices: \n")
	response := fmt.Sprintf(`{"Stock name": "%s", "latest_price": "%s"}`, stockSymbol, latestData.Close)
	fmt.Println("Response:", string(response))
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(response))
	log.Print("Gauge setting init..")
	PriceGauge.WithLabelValues(stockSymbol).Set(result)
	log.Print("Gauge setting finished..")
	if !isChannelAvailable(resultCh) {
		return errors.New("Buffered channel not found")
	}
	log.Print("Channel is available")
	resultCh <- response
	log.Print("Funtion returning")
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

// Helper function
func isChannelAvailable(ch chan string) bool {
	if ch == nil {
		return false
	}

	// Check if the channel is closed
	select {
	case _, ok := <-ch:
		if !ok {
			return false
		}
	default:
	}

	return true
}
