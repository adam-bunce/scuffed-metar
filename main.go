package main

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/serve"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", serve.HandleIndex)
	log.Printf("Server Listening on port: %d\n", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil)
}
