package serve

import (
	"bytes"
	"cmp"
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/stats"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"net/http"
	"slices"
	"strings"
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
	// file, _ := os.ReadFile("serve/index.html") // for local dev
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
		"nl2br": func(str string) template.HTML {
			// just to format taf properly
			hold := strings.ReplaceAll(str, "  ", "&nbsp;&nbsp;")
			hold = strings.ReplaceAll(hold, "\n", "<br>")
			// hold = strings.ReplaceAll(hold, "TEMPO", "<br>TEMPO") TODO format tempo properly
			return template.HTML(hold)
		},
		// check if its CYLJ, CYSF if it is do/don't draw
		"getMetar": func(metars []types.MetarInfo, str string) []string {
			for _, metar := range metars {
				if metar.AirportCode == str {
					return metar.MetarInfo
				}
			}
			return nil
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
	// handlers get put into go routines so lock data incase
	currentData.Lock()
	defer currentData.Unlock()

	var metarData []types.MetarInfo

	// TODO speed test this to see if its actually faster
	results := make(chan []types.MetarInfo, 4)
	go func() { results <- pull.GetAllCamecoData() }()
	go func() { results <- []types.MetarInfo{pull.GetPointsNorthMetar()} }()
	go func() { results <- pull.GetAllHighwayData() }()
	go func() { results <- pull.GetNavCanadaMetars() }()

	for i := 0; i < 4; i++ {
		metarData = append(metarData, <-results...)
	}

	// prevent order from changing in ui
	slices.SortFunc(metarData, func(a, b types.MetarInfo) int {
		return cmp.Compare(strings.ToLower(a.AirportCode), strings.ToLower(b.AirportCode))
	})

	currentData.MetarData = metarData
	currentData.LastUpdate = time.Now().UTC()

	globals.Logger.Println("DONE Updating METAR Data")

	cachedTemplate.Reset()
	err := indexTemplate.Execute(&cachedTemplate, &currentData)
	if err != nil {
		globals.Logger.Printf("Failed to execute index template err: %v\n", err)
		return
	}
}
