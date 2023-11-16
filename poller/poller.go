package poller

// new image to poll and monitor not an api
import (
	"dopemeth/services"
	"log"
	"time"
)

var (
	priceEntry  float64 = -1
	priceLatest float64
	gain        float64 = 0
	loss        float64 = 0
)

// Updates the stock price and measure appropriate gain / loss percentage
func updateState(v services.Output) {
	if priceEntry == -1 {
		log.Print("Price val set")
		priceEntry = v.Price
	}
	if v.Price > priceLatest {
		gain = getgain(priceLatest, v.Price)
		log.Print("gain val got")
		if gain > 1 {
			log.Print("gain val set")
			services.Gaincounter.WithLabelValues(v.Name).Inc()
		}
	} else {
		loss = getloss(priceLatest, v.Price)
		log.Print("loss val got")
		if loss > 1 {
			log.Print("loss val set")
			services.Losscounter.WithLabelValues(v.Name).Inc()
		}
	}
	priceLatest = v.Price
}

func getgain(v1 float64, v2 float64) float64 {
	g := ((v2 - v1) / v1) * 100
	log.Printf("loss %f", g)
	return g
}

func getloss(v1 float64, v2 float64) float64 {
	l := ((v1 - v2) / v1) * 100
	log.Printf("loss %f", l)
	return l
}

func PollApi() error {
	for {
		log.Print("Sending request..")
		val, err := services.FetchAndPrintStockPrices()
		if err != nil {
			return err
		}
		log.Print("Function crossed")
		updateState(val)
		time.Sleep(1 * time.Minute)
		log.Print("Sleep crossed")
	}
}
