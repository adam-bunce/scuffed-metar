package serve

import (
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/stats"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"net/http"
	"os"
	"time"
)

//go:embed index.html
var indexTemplateString string

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateString))
var currentData = types.IndexData{}

// GetTemplate is used when working locally to avoid having to compile for every html update
func GetTemplate() *template.Template {
	file, _ := os.ReadFile("serve/index.html")
	return template.Must(template.New("index").Parse(string(file)))
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	// var indexTemplate = GetTemplate()

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

	err := indexTemplate.Execute(w, &currentData)
	if err != nil {
		globals.Logger.Printf("Failed to execute index template err: %v\n", err)
		return
	}
}

func CatchAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func updateData() {
	globals.Logger.Println("Updating METAR Data")
	currentData.Lock()
	defer currentData.Unlock()

	var metarData []types.MetarInfo
	metarData = append(metarData, pull.GetAllCamecoData()...)
	metarData = append(metarData, pull.GetAllHighwayData()...)
	metarData = append(metarData, pull.GetPointsNorthMetar())

	currentData.MetarData = metarData
	currentData.LastUpdate = time.Now().UTC()

	globals.Logger.Println("DONE Updating METAR Data")
}
