package serve

import (
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
var indexTemplate = LoadTemplate("", "index", indexTemplateString)

//go:embed gfa.html
var gfaTemplateString string
var gfaTemplate = LoadTemplate("", "gfa", gfaTemplateString)

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

	{"", "CYXE", "", nil, types.WeatherInfo{}},
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

var indexData = types.IndexData{
	Version:            globals.Version,
	LastUpdate:         time.Time{},
	AirportInformation: airportInfo,
}

var gfaData = types.GfaPageData{
	Version: globals.Version,
}

func LoadTemplate(templatePath string, templateName string, templateStr ...string) *template.Template {
	var templateString string

	templateFunctions := template.FuncMap{"formatTaf": func(taf string) template.HTML {
		res := strings.Replace(taf, "\n", "<br>", -1)
		res = strings.Replace(res, " ", "&nbsp;", -1)
		return template.HTML(res)
	},
	}

	globals.Logger.Printf("Loading %s Template", templateName)
	if globals.Env == "local" && templatePath != "" {
		globals.Logger.Printf("Getting %s from disk", templatePath)
		file, err := os.ReadFile(templatePath)
		if err != nil {
			globals.Logger.Printf(err.Error())
		}
		templateString = string(file)
	} else if len(templateStr) > 0 {
		globals.Logger.Printf("Loading template from templateStr param")
		return template.Must(template.New(templateName).Funcs(templateFunctions).Parse(templateStr[0]))
	}

	return template.Must(template.New(templateName).Funcs(templateFunctions).Parse(templateString))
}

func TryUpdateGFAData() {
	start := time.Now()
	gfaData.Lock()
	defer gfaData.Unlock()

	if !(time.Since(gfaData.LastUpdate) > 30*time.Second) {
		return
	}

	ids, err := pull.GetGFAImageIds()
	if err != nil {
		gfaData.Error = err
		gfaData.GfaInfo = types.GfaInfo{}
		return
	}

	gfaData.GfaInfo = ids
	gfaData.LastUpdate = time.Now().UTC()

	globals.Logger.Printf("Updated GFA data in %d ms", time.Since(start).Milliseconds())
}

func TryUpdateMETARData() {
	start := time.Now()
	indexData.Lock()
	defer indexData.Unlock()

	if !(time.Since(indexData.LastUpdate) > 30*time.Second) {
		return
	}

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

	indexData.LastUpdate = time.Now().UTC()
	// Update current data w/ new pulled data
	for _, pulledAirport := range metarData {
		for j, currentDataAirports := range indexData.AirportInformation {
			if currentDataAirports.AirportCode == pulledAirport.AirportCode {
				indexData.AirportInformation[j].Metar = pulledAirport.Metar
				indexData.AirportInformation[j].Taf = pulledAirport.Taf
				indexData.AirportInformation[j].Error = pulledAirport.Error
			}
		}
	}

	globals.Logger.Printf("Updated METAR data in %d ms", time.Since(start).Milliseconds())
}
