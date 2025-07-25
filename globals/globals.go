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
var Version = "1.13.8" // for cache busting

var MqttPass = os.Getenv("MQTT_PASS")
var MqttUser = os.Getenv("MQTT_USER")

var VersionHistory = []struct {
	SemVer string
	Change string
}{
	{SemVer: "v" + Version, Change: "cet2 tls error fix"},
	{SemVer: "v1.13.7", Change: "keep stale data when NavCan fails, add warning"},
	{SemVer: "v1.13.6", Change: "add cams for cypa"},
	{SemVer: "v1.13.5", Change: "add south east cam for cyll"},
	{SemVer: "v1.13.4", Change: "add ccb2"},
	{SemVer: "v1.13.3", Change: "add mets option to trip page"},
	{SemVer: "v1.13.2", Change: "add mets"},
	{SemVer: "v1.13.1", Change: "fix buggy trip upper winds checkbox and inconsistent button heights"},
	{SemVer: "v1.13.0", Change: "add WAAS section"},
	{SemVer: "v1.12.4", Change: "add CYYL"},
	{SemVer: "v1.12.3", Change: "bump year, add scuffed meat name"},
	{SemVer: "v1.12.2", Change: "add points north cams"},
	{SemVer: "v1.12.1", Change: "replace CCL3 with CYOD"},
	{SemVer: "v1.12.0", Change: "add winds to trip section"},
	{SemVer: "v1.11.2", Change: "add uranium city METAR"},
	{SemVer: "v1.11.1", Change: "scrape image urls for highways (update if cam is down)"},
	{SemVer: "v1.11.0", Change: "add radius 10nm to notams"},
	{SemVer: "v1.10.0", Change: "navcan fix + pull navcan every 2min"},
	{SemVer: "v1.9.13", Change: "fix cameco metar error"},
	{SemVer: "v1.9.12", Change: "tweak stat internals again"},
	{SemVer: "v1.9.11", Change: "tweak stat internals"},
	{SemVer: "v1.9.10", Change: "add site names"},
	{SemVer: "v1.9.9", Change: "add CYYN, CYXH, CYTH, CYQV. (maybe?) Fix METAR Formatting"},
	{SemVer: "v1.9.8", Change: "add CYLL"},
	{SemVer: "v1.9.7", Change: "add CYQD"},
	{SemVer: "v1.9.6", Change: "fix CET2, add AWOS frequencies"},
	{SemVer: "v1.9.5", Change: "allow CYOD notam in trips section"},
	{SemVer: "v1.9.4", Change: "fix metar spacing"},
	{SemVer: "v1.9.3", Change: "tweak print spacing"},
	{SemVer: "v1.9.2", Change: "tweaked print typography"},
	{SemVer: "v1.9.1", Change: "add instructions"},
	{SemVer: "v1.9.0", Change: "update trip section"},
	{SemVer: "v1.8.0", Change: "add trip section"},
	{SemVer: "v1.7.8", Change: "fix whitescreen, add back load stats"},
	{SemVer: "v1.7.7", Change: "add cams, fix link underlines"},
	{SemVer: "v1.7.6", Change: "add CYMM, CYSM, CYPY"},
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
