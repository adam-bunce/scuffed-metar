package cycles

import (
	"fmt"
	"github.com/adam-bunce/scuffed-metar/globals"
	"sync/atomic"
	"time"
)

var serveCount atomic.Uint64

func IncServeCount() {
	serveCount.Add(1)
}

func StatResetCycle() {
	for {
		now := time.Now().UTC()
		if serveCount.Load() != 0 {
			globals.SendWebhook(fmt.Sprintf(":arrows_clockwise: - METAR/GFA pages loaded %d time(s)", serveCount.Load()))
		}
		serveCount.Store(0)

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
