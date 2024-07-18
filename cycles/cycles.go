package cycles

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"sync/atomic"
	"time"
)

var metarGfaServeCount atomic.Uint64
var tripServeCount atomic.Uint64

func IncMetarGfaServeCount() {
	metarGfaServeCount.Add(1)
}

func IncTripServeCount() {
	tripServeCount.Add(1)
}

func StatResetCycle() {
	for {
		now := time.Now().UTC()
		if metarGfaServeCount.Load() != 0 {
			globals.SendWebhook(fmt.Sprintf(":arrows_clockwise:   __Page Loads__\\n- METAR/GFA pages loaded %d time(s)\\n- TRIP page loaded %d time(s)", metarGfaServeCount.Load(), tripServeCount.Load()))
		}
		metarGfaServeCount.Store(0)
		tripServeCount.Store(0)

		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		sleepDuration := nextMidnight.Sub(now)
		globals.Logger.Printf("StatResetCycle sleeping for %v", sleepDuration)
		time.Sleep(nextMidnight.Sub(now))
	}
}

func EveryTwo(function func()) {
	for {
		function()
		time.Sleep(2 * time.Minute)
	}
}
