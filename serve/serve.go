package serve

import (
	_ "embed"
	"fmt"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/types"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed index.html
var indexTemplateString string

var indexTemplate = template.Must(template.New("index").Parse(indexTemplateString))
var currentData = types.IndexData{}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(fmt.Sprintf("%s %s %s", r.Proto, r.Method, r.RequestURI))

	w.Header().Set("Content-Type", "text/html")

	// update every 30 seconds max
	if currentData.LastUpdate.Before(time.Now().Add(time.Second * -30)) {
		updateData()
	}

	err := indexTemplate.Execute(w, &currentData)
	if err != nil {
		log.Printf("Failed to execute index template err: %v\n", err)
		return
	}
}

func updateData() {
	log.Println("Updating METAR Data")
	currentData.Lock()
	defer currentData.Unlock()

	var metarData []types.MetarInfo
	metarData = append(metarData, pull.GetAllCamecoData()...)
	metarData = append(metarData, pull.GetAllHighwayData()...)
	metarData = append(metarData, pull.GetPointsNorthMetar())

	currentData.MetarData = metarData
	currentData.LastUpdate = time.Now().UTC()
}
