package globals

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var ServerPort = 80
var webhook = os.Getenv("WEBHOOK_URL")
var Env = "prod"
var Version = "1.7.6" // for cache busting

var VersionHistory = []struct {
	SemVer string
	Change string
}{
	{SemVer: "v" + Version, Change: "add CYMM, CYSM, CYPY"},
	{SemVer: "v1.7.5", Change: "fix cameco timeouts"},
	{SemVer: "v1.7.4", Change: "add CYQR, fix setError, update version history internals"},
	{SemVer: "v1.7.3", Change: "fix winds error not resetting"},
	{SemVer: "v1.7.2", Change: "fix buttons jankyness"},
	{SemVer: "v1.7.1", Change: "update cameco, rename upper winds link"},
	{SemVer: "v1.7.0", Change: "added info page & metar src url's"},
	{SemVer: "v1.6.7", Change: "added request timeout"},
	{SemVer: "v1.6.6", Change: "fix hanging requests killing metar page"},
	{SemVer: "v1.6.5", Change: "limited cet2/ccl3 readouts to 5"},
	{SemVer: "v1.6.4", Change: "added cet2/ccl3"},
}

var Timeout = time.Second * 5

var Logger = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds|log.Ldate)

var Client = &http.Client{Timeout: Timeout}

func SendWebhook(message string) {
	_, err := http.Post(webhook, "application/json", strings.NewReader(fmt.Sprintf(`{"content": "%s"}`, message)))
	if err != nil {
		Logger.Printf("Failed to send webhook err: %v\n", err)
	}
}
