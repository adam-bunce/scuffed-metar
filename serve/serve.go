package serve

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed static
var files embed.FS
var fileServer = http.FileServer(http.FS(files))

//go:embed index.html
var indexTemplateString string

var indexTemplate *template.Template = LoadTemplate()

var cachedTemplate bytes.Buffer

var hcUrl = "http://highways.glmobile.com"

func hc(num int) []string {
	var res []string
	for i := 1; i < num+1; i++ {
		res = append(res, fmt.Sprintf("/ptz%s.jpg", strconv.Itoa(i)))
	}
	return res
}

// this is the order displayed in the UI
var airportInfo = []types.AirportInfo{
	{"Cigar Lake", "CJW7", "", nil, types.WeatherInfo{}},
	{"McArthur River", "CKQ8", "", nil, types.WeatherInfo{}},
	{"Collins Bay / Rabbit Lake", "CYKC", "", nil, types.WeatherInfo{}},
	{"Key Lake", "CYKJ", "", nil, types.WeatherInfo{}},

	{"Points North", "CYNL", "", nil, types.WeatherInfo{}},

	{"", "CYPA", "", nil, types.WeatherInfo{}},
	{"", "CYVC", "", nil, types.WeatherInfo{}},

	{"Fond du Lac", "CZFD", hcUrl + "/fonddulac", hc(2), types.WeatherInfo{}},
	{"Wollaston", "CZWL", hcUrl + "/wollaston", hc(2), types.WeatherInfo{}},
	{"Ile A La Crosse", "CJF3", hcUrl + "/ilealacrosse", hc(2), types.WeatherInfo{}},
	{"Cumberland House", "CJT4", hcUrl + "/cumberlandhouse", hc(2), types.WeatherInfo{}},
	{"La Loche", "CJL4", hcUrl + "/laloche", hc(2), types.WeatherInfo{}},
	{"Patuanak", "CKB2", hcUrl + "/patuanak", hc(1), types.WeatherInfo{}},
	{"Pelican Narrows", "CJW4", hcUrl + "/pelican", hc(1), types.WeatherInfo{}},
	{"Pinehouse", "CZPO", hcUrl + "/pinehouse", hc(2), types.WeatherInfo{}},
	{"Buffalo Narrows", "CYVT", hcUrl + "/buffalonarrows", hc(2), types.WeatherInfo{}},
	{"Hudson Bay", "CYHB", hcUrl + "/hudsonbay", hc(2), types.WeatherInfo{}},
	{"Stony Rapids", "CYSF", hcUrl + "/stonyrapids", hc(2), types.WeatherInfo{}},
	{"Sandy Bay", "CJY4", hcUrl + "/sandybay", []string{"/ptz.jpg", "/ptz2.jpg", "/ptz3.jpg"}, types.WeatherInfo{}},
	{"Meadow Lake", "CYLJ", hcUrl + "/meadowlake", hc(2), types.WeatherInfo{}},
	{"Uranium City", "CYBE", hcUrl + "/uranium", hc(1), types.WeatherInfo{}},

	{"Charlot River", "CJP9", "http://saskpower.glmobile.com/charlot", []string{"/runway.jpg", "/hill.jpg"}, types.WeatherInfo{}},
}

var currentData = types.IndexData{
	Version:            globals.Version,
	LastUpdate:         time.Time{},
	AirportInformation: airportInfo,
}

func LoadTemplate() *template.Template {
	globals.Logger.Printf("Loading Template")
	if globals.Env == "local" {
		globals.Logger.Printf("Getting index.html from disk")
		file, err := os.ReadFile("serve/index.html")
		if err != nil {
			globals.Logger.Printf(err.Error())
		}
		indexTemplateString = string(file)
	}

	templateFunctions := template.FuncMap{
		"formatTaf": func(taf string) template.HTML {
			res := strings.Replace(taf, "\n", "<br>", -1)
			res = strings.Replace(res, " ", "&nbsp;", -1)
			return template.HTML(res)
		},
	}

	return template.Must(template.New("index").Funcs(templateFunctions).Parse(indexTemplateString))
}

func UpdateData() {
	start := time.Now()
	currentData.Lock()
	defer currentData.Unlock()

	dataChan := make(chan types.WeatherPullInfo)

	var wg sync.WaitGroup
	wg.Add(4)
	go pull.GetAllCamecoData(dataChan, &wg)
	go pull.GetAllHighwayData(dataChan, &wg)
	go pull.GetPointsNorthMetar(dataChan, &wg)
	go pull.GetNavCanadaMetars(dataChan, &wg)

	// close chan so read loop doesn't hang
	go func() {
		wg.Wait()
		close(dataChan)
	}()

	var metarData []types.WeatherPullInfo
	for metar := range dataChan {
		metarData = append(metarData, metar)
	}

	currentData.LastUpdate = time.Now().UTC()
	// Update current data w/ new pulled data
	for _, pulledAirport := range metarData {
		for j, currentDataAirports := range currentData.AirportInformation {
			if currentDataAirports.AirportCode == pulledAirport.AirportCode {
				currentData.AirportInformation[j].Metar = pulledAirport.Metar
				currentData.AirportInformation[j].Taf = pulledAirport.Taf
			}
		}
	}

	end := time.Now()
	globals.Logger.Printf("Updated METAR data in %s", end.Sub(start))
}
