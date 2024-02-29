package serve

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/cycles"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/types"
	"net/http"
	"slices"
	"strings"
	"time"
)

func HandleAll(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	cycles.IncServeCount()

	if globals.Env == "local" {
		indexTemplate = LoadTemplate("serve/pages/index.html", "index")
	}

	if globals.Env != "local" {
		TryUpdateMETARData()
	}

	indexTemplate.Execute(w, &indexData)
}

func HandleGfa(w http.ResponseWriter, r *http.Request) {
	cycles.IncServeCount()
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		gfaTemplate = LoadTemplate("serve/pages/gfa.html", "gfa")
	}

	if globals.Env != "local" {
		TryUpdateGFAData()
	}

	gfaTemplate.Execute(w, &gfaData)
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	w.Header().Set("Cache-Control", "max-age=604800")

	if globals.Env == "local" {
		globals.Logger.Printf("using http.Dir NOT embedded fs")
		http.FileServer(http.Dir("serve/")).ServeHTTP(w, r)
	} else {
		fileServer.ServeHTTP(w, r)
	}
}

func HandleNotam(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	airportCodes := r.URL.Query()["airport"]

	if globals.Env == "local" {
		notamTemplate = LoadTemplate("serve/pages/notam.html", "notam")
	}

	var data []types.NotamData
	if len(airportCodes) > 0 {
		data, _ = pull.GetNotam(airportCodes)
	}
	notamData.LastUpdate = time.Now().UTC()

	notamTemplate.Execute(w, map[string]interface{}{"notam": &notamData, "data": data})
}

func HandleWinds(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	airportCodes := r.URL.Query()["airport"]

	if globals.Env == "local" {
		windsTemplate = LoadTemplate("serve/pages/winds.html", "winds")
	}

	var data []types.WindsData
	if len(airportCodes) > 0 {
		data, windsData.Error = pull.GetWinds(airportCodes)
	}
	windsData.LastUpdate = time.Now().UTC()
	windsTemplate.Execute(w, map[string]interface{}{"winds": &windsData, "data": data})
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if r.Method == http.MethodPost {
		msg := r.PostFormValue("info")
		if strings.TrimSpace(msg) != "" {
			globals.SendWebhook(fmt.Sprintf(":pencil: - %s", msg))
		}
	}

	if globals.Env == "local" {
		infoTemplate = LoadTemplate("serve/pages/info.html", "info")
	}

	infoTemplate.Execute(w, map[string]interface{}{"Version": globals.Version, "VersionHistory": globals.VersionHistory})
}

func HandleTrip(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	airportCodes := r.URL.Query()["airport"]

	if globals.Env != "local" {
		TryUpdateMETARData()
	}

	if globals.Env == "local" {
		tripTemplate = LoadTemplate("serve/pages/trip.html", "trip")
	}

	var selectedAirportsInfo []types.AirportInfo
	for _, airport := range indexData.AirportInformation {
		if slices.Contains(airportCodes, airport.AirportCode) {
			selectedAirportsInfo = append(selectedAirportsInfo, airport)
		}
	}

	// NOTE: notam data has all airports with notams which matches the airports with metars
	err := tripTemplate.Execute(w, map[string]any{"airportInfo": selectedAirportsInfo, "options": notamData, "LastUpdate": &indexData.LastUpdate})
	if err != nil {
		panic(err)
		return
	}
}
