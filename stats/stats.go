package stats

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
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

		globals.SendWebhook(fmt.Sprintf(":arrows_clockwise: - METAR/GFA pages loaded %d time(s)", serveCount.Load()))
		serveCount.Store(0)

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		sleepDuration := nextMidnight.Sub(now)
		globals.Logger.Printf("StatResetCycle sleeping for %v", sleepDuration)
		time.Sleep(nextMidnight.Sub(now))
	}
}
