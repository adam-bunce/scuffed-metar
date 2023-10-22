package main

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/serve"
	"github.com/adam-bunce/scuffed-metar/stats"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", serve.HandleIndex)

	go stats.StatResetCycle()

	log.Printf("Server Listening on port: %d\n", globals.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", globals.ServerPort), nil)
}
