package main

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/cycles"
	"github.com/adam-bunce/scuffed-metar/globals"
	"github.com/adam-bunce/scuffed-metar/serve"
	"net/http"
)

func main() {
	globals.SendWebhook(fmt.Sprintf("Deployed %s: %s",
		globals.VersionHistory[0].SemVer,
		globals.VersionHistory[0].Change))

	serve.TryUpdateGFAData()
	serve.TryUpdateMETARData()

	go cycles.EveryTwo(serve.TryNavCanUpdate)
	go cycles.EveryTwo(serve.TryCamecoUpdate)
	go cycles.StatResetCycle()

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", serve.HandleStatic)
	mux.HandleFunc("/graphic-area-forecast", serve.HandleGfa)
	mux.HandleFunc("/notam", serve.HandleNotam)
	mux.HandleFunc("/winds", serve.HandleWinds)
	mux.HandleFunc("/info", serve.HandleInfo)
	mux.HandleFunc("/trip", serve.HandleTrip)
	mux.HandleFunc("/waas", serve.HandleWaas)
	mux.HandleFunc("/mets", serve.HandleMets)
	mux.HandleFunc("/", serve.HandleAll)

	http.ListenAndServe(fmt.Sprintf(":%d", globals.ServerPort), mux)
}
