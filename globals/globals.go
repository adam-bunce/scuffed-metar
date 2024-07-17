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
var Version = "1.9.11" // for cache busting

var MqttPass = os.Getenv("MQTT_PASS")
var MqttUser = os.Getenv("MQTT_USER")

var VersionHistory = []struct {
	SemVer string
	Change string
}{
	{SemVer: "v" + Version, Change: "tweak stat internals"},
	{SemVer: "v1.9.10" + Version, Change: "add site names"},
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
