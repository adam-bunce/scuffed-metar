package serve

import (
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/stats"
	"net/http"
	"time"
)

func HandleAll(w http.ResponseWriter, r *http.Request) {
	globals.Logger.Printf("%s %s %s", r.Proto, r.Method, r.RequestURI)

	if r.Method != http.MethodGet || r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if globals.Env == "local" {
		indexTemplate = LoadTemplate()
	}

	stats.IncServeCount()

	// update every 30 seconds max
	if currentData.LastUpdate.Before(time.Now().Add(time.Second*-30)) && globals.Env != "local" {
		UpdateData()
	}

	cachedTemplate.Reset()
	err := indexTemplate.Execute(&cachedTemplate, &currentData)
	if err != nil {
		globals.Logger.Printf(err.Error())
	}

	w.Write(cachedTemplate.Bytes())
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