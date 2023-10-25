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
	// TODO write 404 page
	http.HandleFunc("/", serve.HandleIndex)

	go stats.StatResetCycle()

	globals.Logger.Printf("Server Listening on port: %d\n", globals.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ServerPort), nil)
}
