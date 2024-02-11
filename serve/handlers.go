package serve

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/stats"
	"github.com/adam-bunce/scuffed-metar/types"
	"net/http"
	"strings"
	"time"
)

func HandleAll(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if globals.Env == "local" {
		indexTemplate = LoadTemplate("serve/pages/index.html", "index")
	}

	stats.IncServeCount()

	if globals.Env != "local" {
		TryUpdateMETARData()
	}

	indexTemplate.Execute(w, &indexData)
}

func HandleGfa(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		gfaTemplate = LoadTemplate("serve/pages/gfa.html", "gfa")
	}

	stats.IncServeCount()

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
	var err error
	if len(airportCodes) > 0 {
		data, err = pull.GetWinds(airportCodes)
		if err != nil {
			windsData.Error = err
		}
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

	infoTemplate.Execute(w, map[string]interface{}{"Version": globals.Version})
}
