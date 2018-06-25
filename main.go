package main


import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"log"
	"strconv"
	"strings"
	"net/http"
	"encoding/csv"
	"io"
	"os"
)

type Statement struct {
	credit int `json:"credit"`
	debit  int `json:"debit"`
}

func parseCSV(r io.Reader)(Statement, error) {

	csvFile := csv.NewReader(r)

	records, err := csvFile.ReadAll()
	if err != nil {
		log.Println(err)
	}

	// skips first line of csv
	records = records[1:]
	credit := 0
	debit := 0

	for _, record := range records {
		amount, err := strconv.Atoi(record[3])
		if err != nil {
			log.Println(err)
		}

		if record[2] == "C" {
			credit += amount
		} else {
			debit += amount
		}
	}
	
	s := Statement{credit, debit}

	return s, err

}


func drawChart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	// makes sure type is of multipart-form
	if content := r.Header.Get("Content-Type"); !strings.HasPrefix(content, "multipart/form-data") {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}


	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := parseCSV(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(s)

	sbc := chart.BarChart{
		Title:      "Test Bar Chart",
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: []chart.Value{
			{Value: float64(s.credit), Label: "Blue"},
			{Value: float64(s.debit), Label: "Green"},
		},
	}

	w.Header().Set("Content-Type", "image/png")
	err = sbc.Render(chart.PNG, w)
	if err != nil {
		log.Printf("Error rendering chart: %v\n", err)
		return
	}
}


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/draw", drawChart)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	port := 3000
	if postStr, ok := os.LookupEnv("PORT"); ok {
		if i, err := strconv.Atoi(postStr); err == nil {
			port = i
		}
	}

	log.Printf("starting server on http://localhost:%d/\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}