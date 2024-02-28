package cycles

import (
	"github.com/adam-bunce/scuffed-metar/serve"
	"time"
)

func CamecoPullCycle() {
	for {
		serve.TryCamecoUpdate()
		time.Sleep(2 * time.Minute)
	}
}
