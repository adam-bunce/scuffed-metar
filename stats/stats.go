package stats

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"log"
	"sync/atomic"
	"time"
)

var serveCount atomic.Int64 // can have multiple requests

func IncServeCount() {
	serveCount.Add(1)
}

func StatResetCycle() {
	for {
		now := time.Now().UTC()

		globals.SendWebhook(fmt.Sprintf("%s - %d", now.Format("01-02-2006"), serveCount.Load()))
		serveCount.Store(0)

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		sleepDuration := nextMidnight.Sub(now)
		log.Printf("StatResetCycle sleeping for %v", sleepDuration)
		time.Sleep(nextMidnight.Sub(now))
	}
}
