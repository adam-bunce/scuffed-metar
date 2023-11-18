package serve

import (
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/stats"
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
		indexTemplate = LoadTemplate("serve/index.html", "index")
	}

	stats.IncServeCount()

	// update every 30 seconds max
	if indexData.LastUpdate.Before(time.Now().Add(time.Second*-30)) && globals.Env != "local" {
		UpdateIndexData()
	}

	// this looks like an unnecesscary execute? cant we just send the bytes?
	cachedIndexTemplate.Reset()
	err := indexTemplate.Execute(&cachedIndexTemplate, &indexData)
	if err != nil {
		globals.Logger.Printf(err.Error())
	}

	w.Write(cachedIndexTemplate.Bytes())
}

func HandleGfa(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		gfaTemplate = LoadTemplate("serve/gfa.html", "gfa")
	}

	stats.IncServeCount()

	// update every 30 seconds max
	if gfaData.LastUpdate.Before(time.Now().Add(time.Second*-30)) && globals.Env != "local" {
		UpdateGfaData()
	}

	// update template
	cachedGfaTemplate.Reset()
	err := gfaTemplate.Execute(&cachedGfaTemplate, &gfaData)
	if err != nil {
		globals.Logger.Printf(err.Error())
	}

	w.Write(cachedGfaTemplate.Bytes())
}

func HandleGfaSubRoute(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	URIParts := strings.Split(r.RequestURI, "/")
	targetGFA := URIParts[len(URIParts)-1]

	if gfaData.LastUpdate.Before(time.Now().Add(time.Second*-30)) && globals.Env != "local" {
		UpdateGfaData()
	}

	executeData := struct {
		Type          string
		Version       string
		Id            string
		StartValidity string
		LastUpdate    time.Time
	}{}
	// find matching data to send
	if strings.Contains(targetGFA, "CLDWX") {
		for _, entry := range gfaData.CloudsWeather {
			if strings.Replace(targetGFA, "CLDWX-", "", -1) == entry.StartValidity {
				executeData.Type = "CLDWX"
				executeData.Id = entry.Id
				executeData.StartValidity = entry.StartValidity
				executeData.Version = globals.Version
				executeData.LastUpdate = gfaData.LastUpdate
			}
		}
	} else {
		for _, entry := range gfaData.IcingTurbFreezing {
			if strings.Replace(targetGFA, "TURBC-", "", -1) == entry.StartValidity {
				executeData.Type = "TURBC"
				executeData.Id = entry.Id
				executeData.StartValidity = entry.StartValidity
				executeData.Version = globals.Version
				executeData.LastUpdate = gfaData.LastUpdate
			}
		}
	}

	// send stuff
	err := gfaSubRouteTemplate.Execute(w, &executeData)
	if err != nil {
		globals.Logger.Printf(err.Error())
	}

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

// TODO logging with the status code
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// TODO add status code
		globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)
	})
}
