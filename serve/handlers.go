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

	cycles.IncMetarGfaServeCount()

	if globals.Env == "local" {
		indexTemplate = LoadTemplate("serve/pages/index.html", "index")
	}

	if globals.Env != "local" {
		TryUpdateMETARData()
	}

	indexTemplate.Execute(w, &indexData)
}

func HandleGfa(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)
	cycles.IncMetarGfaServeCount()

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
		data = pull.GetAllNotams(airportCodes)
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
	cycles.IncTripServeCount()

	// airports that we can pull data from CFPS
	validAirports := []string{"CJW7", "CKQ8", "CYKC", "CYKJ", "CYNL", "CYXE",
		"CYPA", "CYVC", "CZFD", "CZWL", "CJF3", "CJT4",
		"CJL4", "CKB2", "CJW4", "CZPO", "CYVT", "CYHB",
		"CYSF", "CJY4", "CYLJ", "CYBE", "CJP9", "CYFO", "CYQW",
		"CET2", "CYOD", "CYQR",
		"CYMM", "CYSM", "CYPY",
		"CYOD", "CYQD", "CYLL",
		"CYYN", "CYXH", "CYTH", "CYQV", "CYYL", "CCB2",
	}

	airportCodes := r.URL.Query()["airport"]
	includeWinds := r.URL.Query()["winds"]
	includeMets := r.URL.Query()["mets"]

	// remove sites sent by user that don't exist
	var filteredValidAirports []string
	for _, code := range airportCodes {
		if slices.Contains(validAirports, code) {
			filteredValidAirports = append(filteredValidAirports, code)
		}
	}

	// ik this is bad fyi
	var notamActualData []types.NotamData
	var metsData map[string][]string
	var data []types.WindsData
	if globals.Env != "local" && len(filteredValidAirports) > 0 {
		TryUpdateMETARData()
		TryUpdateGFAData()
		notamActualData = pull.GetAllNotams(filteredValidAirports)

		if len(includeWinds) > 0 && includeWinds[0] == "true" {
			data, windsData.Error = pull.GetWinds(windsData.WindsOptions)
		}

		if len(includeMets) > 0 && includeMets[0] == "true" {
			metsData = pull.GetAllSigmetAirmet(pull.NavCanSites)
		}
	}

	if globals.Env == "local" {
		tripTemplate = LoadTemplate("serve/pages/trip.html", "trip")
	}

	var selectedAirportsMetar []types.AirportInfo
	for _, airport := range indexData.AirportInformation {
		if slices.Contains(filteredValidAirports, airport.AirportCode) {
			selectedAirportsMetar = append(selectedAirportsMetar, airport)
		}
	}

	// NOTE: notam data has all airports with notams which matches the airports with metars
	tripTemplate.Execute(w, map[string]any{
		"airportInfo": selectedAirportsMetar,
		"options":     notamData,
		"RequestedAt": time.Now().UTC(),
		"gfa":         &gfaData,
		"notam":       &notamActualData,

		// YIKES!
		"data":  &data,
		"winds": &windsData,
		"mets":  &metsData,
	})

}

func HandleWaas(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		waasTemplate = LoadTemplate("serve/pages/waas.html", "info")
	}

	waasTemplate.Execute(w, map[string]interface{}{
		"LastUpdate": time.Now().UTC(),
		"Version":    globals.Version,
	})
}

func HandleMets(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		metsTemplate = LoadTemplate("serve/pages/mets.html", "met")
	}

	data := pull.GetAllSigmetAirmet(append(pull.NavCanSites, "CCB2"))

	metsTemplate.Execute(w, map[string]interface{}{
		"LastUpdate": time.Now().UTC(),
		"Version":    globals.Version,
		"Data":       data,
	})
}
