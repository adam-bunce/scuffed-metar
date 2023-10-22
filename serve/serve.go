package serve

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/stats"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"net/http"
	"time"
)

//go:embed index.html
var indexTemplateString string
var indexTemplate = GetTemplate()

var cachedTemplate bytes.Buffer

var currentData = types.IndexData{
	Cameras: wxCamInfo,
}

var wxCamInfo = []types.WxCam{
	{"ilealacrosse", "CJF3", 2},
	{"cumberlandhouse", "CJT4", 2},
	{"laloche", "CJL4", 2},
	{"patuanak", "CKB2", 1},
	{"pelican", "CJW4", 1},
	{"pinehouse", "CZPO", 2},
	{"buffalonarrows", "CYVT", 2},
	{"hudsonbay", "CYHB", 2},
	{"stonyrapids", "CYSF", 2},
	{"sandybay", "CJY4", 3},
	{"meadowlake", "CYLJ", 2},
	{"uranium", "CYBE", 1},
}

func GetTemplate() *template.Template {
	// file, _ := os.ReadFile("serve/index.html") for local dev
	tmplFuncs := template.FuncMap{
		// makes a slice given a number to iterate over with range
		"makeSlice": func(start, stop int) []int {
			var result []int
			// base 1 intentionally
			for i := start; i < stop+1; i++ {
				result = append(result, i)
			}
			return result
		},
	}
	return template.Must(template.New("index").Funcs(tmplFuncs).Parse(indexTemplateString))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Println(fmt.Sprintf("%s %s %s", r.Proto, r.Method, r.RequestURI))

	if r.URL.Path != "/" {
		CatchAll(w, r)
		return
	}
	stats.IncServeCount()
	w.Header().Set("Content-Type", "text/html")

	// update every 30 seconds max
	if currentData.LastUpdate.Before(time.Now().Add(time.Second * -30)) {
		updateData()
	}

	w.Write(cachedTemplate.Bytes())
}

func CatchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func updateData() {
	// var indexTemplate = GetTemplate()

	globals.Logger.Println("Updating METAR Data")
	currentData.Lock()
	defer currentData.Unlock()

	var metarData []types.MetarInfo
	metarData = append(metarData, pull.GetAllCamecoData()...)
	metarData = append(metarData, pull.GetPointsNorthMetar())
	metarData = append(metarData, pull.GetAllHighwayData()...)

	currentData.MetarData = metarData
	currentData.LastUpdate = time.Now().UTC()

	globals.Logger.Println("DONE Updating METAR Data")

	err := indexTemplate.Execute(&cachedTemplate, &currentData)
	if err != nil {
		globals.Logger.Printf("Failed to execute index template err: %v\n", err)
		return
	}
}
