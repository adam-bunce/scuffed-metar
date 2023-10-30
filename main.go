package main

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/serve"
	"github.com/adam-bunce/scuffed-metar/stats"
	"net/http"
)

func main() {
	// TODO use new go routing to match requests properly to the methods
	// TODO write logging middleware
	serve.UpdateData()
	go stats.StatResetCycle()

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", serve.HandleStatic)
	mux.HandleFunc("/", serve.HandleAll)

	http.ListenAndServe(fmt.Sprintf(":%d", globals.ServerPort), mux)
}
