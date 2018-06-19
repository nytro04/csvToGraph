package main
 
import (
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"strconv"
	"encoding/csv"
	"log"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "os"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/vgsvg"
	"math"
	"gonum.org/v1/plot/vg/draw"
	"bytes"
)

func renderSVG(p *plot.Plot) string {
	size := 10 * vg.Centimeter
	canvas := vgsvg.New(size, size/vg.Length(math.Phi))
	p.Draw(draw.New(canvas))
	out := new(bytes.Buffer)
	_, err := canvas.WriteTo(out)
	if err != nil {
		panic(err)
	}
	return string(out.Bytes())
}

// GetCsv gets the and returns a processed picture
func GetCsv(filepath string) (string, error) {

		// Open the csv file
		file, err := os.Open(filepath)
		if err != nil {
			log.Println("Could not open csv file", err)
		}
	
		// Read from the file
		reader := csv.NewReader(file)
	
		defer file.Close()
	
		reader.FieldsPerRecord = -1
		credit := 0
		debit := 0
	
		for {
	
			records, err := reader.Read()
			if err == io.EOF {
				break
			}
	
			if records[3] == "AMOUNT" {
				continue
			}
	
			if err != nil {
				log.Printf("Can not read CSV file %s", err)
				continue
			}
	
			price, err := strconv.Atoi(records[3])
			if err != nil {
				log.Printf("Can not retrieve amount %s", err)
			} 
	
			if records[2] == "C" {
				credit += price
			} else {
				debit += price
			}
	
			// fmt.Printf("This is the credit transaction %d", credit)
			// fmt.Printf("This is the debit transaction %d",debit)
			
		}
	
		// fmt.Println("This is the credit transaction ", credit)
		// fmt.Println("This is the debit transaction ",debit)
	
		creditfloat := float64(credit)
		debitfloat  := float64(debit)
	
		// Plotting with gonums
		creditTran := plotter.Values{creditfloat}
		debitTran  := plotter.Values{debitfloat}
	
		p, err := plot.New()
		if err != nil {
			log.Println(err)
		}
		p.Title.Text = "Bar Charts"
		p.Y.Label.Text = "Amount"
		p.X.Label.Text = "Transaction Type"
	
		w := vg.Points(20)
	
		janBar, err := plotter.NewBarChart(creditTran, w)
		if err != nil {
			log.Println(err)
		}
		janBar.LineStyle.Width = vg.Length(0)
		janBar.Color = plotutil.Color(0)
		janBar.Offset = -w
	
		janDeb, err := plotter.NewBarChart(debitTran, w)
		if err != nil {
			log.Println(err)
		}
		janDeb.LineStyle.Width = vg.Length(0)
		janDeb.Color = plotutil.Color(1)

		p.Add(janBar, janDeb)
		p.Legend.Add("Credit", janBar)
		p.Legend.Add("Debit", janDeb)
	
		p.Legend.Top = true
		p.NominalX("One", "Two")

		if err := p.Save(5*vg.Inch, 3*vg.Inch, "janTran.png"); err != nil {
			log.Println(err)
		}
	
	return renderSVG(p), err
}


func upload(w http.ResponseWriter, r *http.Request) {
 
    if r.Method == "GET" {
        // GET
        t, _ := template.ParseFiles("index.html")
 
        t.Execute(w, nil)
 
    } else if r.Method == "POST" {
        // Post
        _, err := GetCsv("work.csv")
        if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
        }
 
 
    } else {
        fmt.Println("Unknown HTTP " + r.Method + "  Method")
    }
}

func defaultCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	x, _ := GetCsv("work.csv")
	fmt.Fprint(w, x)
}

func main() {
    http.HandleFunc("/upload", upload)
	http.HandleFunc("/work.csv", defaultCSV)
    http.ListenAndServe(":8000", nil) // setting listening port
}