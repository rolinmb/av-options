package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

const (
    AVURL1 = "https://www.alphavantage.co/query?function=HISTORICAL_OPTIONS&symbol="
    AVURL2 = "&date="
    AVKEY = ""
)

// Define the struct to match JSON structure
type OptionResponse struct {
	Endpoint string   `json:"endpoint"`
	Message  string   `json:"message"`
	Data     []Option `json:"data"`
}

type Option struct {
	ContractID       string  `json:"contractID"`
	Symbol           string  `json:"symbol"`
	Expiration       string  `json:"expiration"`
	Strike           string  `json:"strike"`
	Type             string  `json:"type"`
	Last             string  `json:"last"`
	Mark             string  `json:"mark"`
	Bid              string  `json:"bid"`
	BidSize          string  `json:"bid_size"`
	Ask              string  `json:"ask"`
	AskSize          string  `json:"ask_size"`
	Volume           string  `json:"volume"`
	OpenInterest     string  `json:"open_interest"`
	Date             string  `json:"date"`
	ImpliedVol       string  `json:"implied_volatility"`
	Delta            string  `json:"delta"`
	Gamma            string  `json:"gamma"`
	Theta            string  `json:"theta"`
	Vega             string  `json:"vega"`
	Rho              string  `json:"rho"`
}

func main() {
    ticker := "SPY"
	url := AVURL1 + ticker + "&apikey=" + AVKEY

	// Fetch JSON
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}

	// Parse JSON
	var options OptionResponse
	if err := json.Unmarshal(body, &options); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(options.Data) == 0 {
		log.Println("no option data returned")
		return
	}

	// Create CSV file
	csvName := "data/" + ticker + "options.csv"
	file, err := os.Create(csvName)
	if err != nil {
		log.Fatalf("failed to create csv file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Use reflection to automatically generate headers from struct field tags
	first := options.Data[0]
	val := reflect.ValueOf(first)
	typ := val.Type()

	headers := []string{}
	for i := 0; i < typ.NumField(); i++ {
		headers = append(headers, typ.Field(i).Tag.Get("json"))
	}
	writer.Write(headers)

	// Write data rows
	for _, opt := range options.Data {
		row := []string{
			opt.ContractID,
			opt.Symbol,
			opt.Expiration,
			opt.Strike,
			opt.Type,
			opt.Last,
			opt.Mark,
			opt.Bid,
			opt.BidSize,
			opt.Ask,
			opt.AskSize,
			opt.Volume,
			opt.OpenInterest,
			opt.Date,
			opt.ImpliedVol,
			opt.Delta,
			opt.Gamma,
			opt.Theta,
			opt.Vega,
			opt.Rho,
		}
		writer.Write(row)
	}

	fmt.Printf("main.go :: %s CSV file %s written successfully.", ticker, csvName)
}
