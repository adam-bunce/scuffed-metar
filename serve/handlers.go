package serve

import (
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/pull"
	"github.com/adam-bunce/scuffed-metar/stats"
	"net/http"
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

	if globals.Env != "local" {
		TryUpdateMETARData()
	}

	indexTemplate.Execute(w, &indexData)
}

func HandleGfa(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if globals.Env == "local" {
		gfaTemplate = LoadTemplate("serve/gfa.html", "gfa")
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
		notamTemplate = LoadTemplate("serve/notam.html", "notam")
	}

	data, _ := pull.GetNotam(airportCodes)

	notamTemplate.Execute(w, map[string]interface{}{"notam": &notamData, "data": data})
}

// TODO logging with the status code
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// TODO add status code
		globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)
	})
}
