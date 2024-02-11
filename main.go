package main

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/serve"
	"github.com/adam-bunce/scuffed-metar/stats"
	"net/http"
)

func main() {
	serve.TryUpdateGFAData()
	serve.TryUpdateMETARData()

	go stats.StatResetCycle()

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", serve.HandleStatic)
	mux.HandleFunc("/graphic-area-forecast", serve.HandleGfa)
	mux.HandleFunc("/notam", serve.HandleNotam)
	mux.HandleFunc("/winds", serve.HandleWinds)
	mux.HandleFunc("/info", serve.HandleInfo)
	mux.HandleFunc("/", serve.HandleAll)

	http.ListenAndServe(fmt.Sprintf(":%d", globals.ServerPort), mux)
}
