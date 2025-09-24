package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"io"
	"net/http"
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
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
    	log.Fatalf("failed to create data directory: %v", err)
	}
	if err := os.MkdirAll("img", os.ModePerm); err != nil {
		log.Fatalf("failed to create img directory: %v", err)
	}
	//currentOptionChain("qqq")
	//historicalOptionChain("qqq", "2025-09-22")
	plotIVSurface("data/qqqchain.csv")
	plotIVSurface("data/qqq2025-09-22chain.csv")
	plotIVSurface("data/qqq2025-09-19chain.csv")
	plotIVSurface("data/qqq2025-09-18chain.csv")
	plotIVSurface("data/qqq2025-09-17chain.csv")
	plotIVSurface("data/qqq2025-09-16chain.csv")
}

func historicalOptionChain(ticker,date string) {
	url := AVURL1 + ticker + AVURL2 + date + "&apikey=" + AVKEY
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
	csvName := "data/" + ticker + date + "chain.csv"
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

	fmt.Printf("main.go :: CSV file %s written successfully.\n", csvName)
}

func currentOptionChain(ticker string) {
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
	csvName := "data/" + ticker + "chain.csv"
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

	fmt.Printf("main.go :: CSV file %s written successfully.\n", csvName)
}

// plotIVSurface reads a CSV and generates separate 3D IV surface plots for calls and puts
func plotIVSurface(csvPath string) {
	// ensure img directory exists
	if err := os.MkdirAll("img", os.ModePerm); err != nil {
		log.Fatalf("plotIVSurface: failed to create img directory: %v", err)
	}

	// open CSV
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("plotIVSurface: open csv: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	// --- read header ---
	header, err := r.Read()
	if err != nil {
		log.Fatalf("plotIVSurface: read header: %v", err)
	}
	indices := map[string]int{}
	for i, h := range header {
		indices[strings.TrimSpace(strings.ToLower(h))] = i
	}
	get := func(name string) int {
		if idx, ok := indices[strings.ToLower(name)]; ok {
			return idx
		}
		log.Fatalf("plotIVSurface: missing column %s", name)
		return -1
	}

	idxStrike := get("strike")
	idxExpiration := get("expiration")
	idxDate := get("date")
	idxIV := get("implied_volatility")
	idxType := get("type") // call/put

	// --- points files ---
	callPoints := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_calls_ivpoints.dat")
	putPoints := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_puts_ivpoints.dat")

	callF, _ := os.Create(callPoints)
	defer callF.Close()
	putF, _ := os.Create(putPoints)
	defer putF.Close()

	callCount, putCount := 0, 0

	// --- process records ---
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}

		typ := strings.ToLower(strings.TrimSpace(rec[idxType]))
		strike, err := strconv.ParseFloat(rec[idxStrike], 64)
		if err != nil {
			continue
		}
		ivStr := strings.TrimSuffix(rec[idxIV], "%")
		iv, err := strconv.ParseFloat(ivStr, 64)
		if err != nil {
			continue
		}
		if strings.Contains(rec[idxIV], "%") {
			iv /= 100.0
		}

		expT, _ := time.Parse("2006-01-02", rec[idxExpiration])
		dateT, _ := time.Parse("2006-01-02", rec[idxDate])
		if expT.IsZero() || dateT.IsZero() || !expT.After(dateT) {
			continue
		}
		T := expT.Sub(dateT).Hours() / (24.0 * 365.0)
		if T <= 0 {
			continue
		}

		line := fmt.Sprintf("%g %g %g\n", strike, T, iv)
		if typ == "call" {
			callF.WriteString(line)
			callCount++
		} else if typ == "put" {
			putF.WriteString(line)
			putCount++
		}
	}

	if callCount == 0 && putCount == 0 {
		log.Println("plotIVSurface: no valid points extracted")
		return
	}

	fmt.Printf("plotIVSurface: wrote %d call points to %s\n", callCount, callPoints)
	fmt.Printf("plotIVSurface: wrote %d put points to %s\n", putCount, putPoints)

	// --- gnuplot script generator ---
	makeScript := func(pointsFile, outPNG, title string) string {
    	return fmt.Sprintf(`
set terminal pngcairo size 1200,800 enhanced font 'Arial,12'
set output '%s'
set title '%s'
set xlabel 'Strike'
set ylabel 'Time to Expiration (Years)'
set zlabel 'Implied Volatility'
set grid
set view 60,30
set dgrid3d 50,50 gauss 2
set pm3d at s
set palette defined (0 'blue', 1 'green', 2 'yellow', 3 'red')
splot '%s' using 1:2:3 with pm3d notitle
`, filepath.ToSlash(outPNG), title, filepath.ToSlash(pointsFile))
	}

	runGnuplot := func(scriptFile string) {
		cmd := exec.Command("gnuplot", scriptFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("plotIVSurface: gnuplot failed: %v", err)
		}
	}

	// --- generate plots ---
	if callCount > 0 {
		callPNG := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_call_ivsurface.png")
		callScript := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_call_ivsurface.gp")
		os.WriteFile(callScript, []byte(makeScript(callPoints, callPNG, "Call IV Surface")), 0644)
		runGnuplot(callScript)
		fmt.Printf("plotIVSurface: call PNG saved to %s\n", callPNG)

		os.Remove(callScript)
		os.Remove(callPoints)
	}

	if putCount > 0 {
		putPNG := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_put_ivsurface.png")
		putScript := filepath.Join("img", strings.TrimSuffix(filepath.Base(csvPath), ".csv")+"_put_ivsurface.gp")
		os.WriteFile(putScript, []byte(makeScript(putPoints, putPNG, "Put IV Surface")), 0644)
		runGnuplot(putScript)
		fmt.Printf("plotIVSurface: put PNG saved to %s\n", putPNG)

		os.Remove(putScript)
		os.Remove(putPoints)
	}
}
