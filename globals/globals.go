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
var Version = "1.7.0" // for cache busting

var Logger = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds|log.Ldate)

var Client = &http.Client{Timeout: time.Second * 3} //

func SendWebhook(message string) {
	_, err := http.Post(webhook, "application/json", strings.NewReader(fmt.Sprintf(`{"content": "%s"}`, message)))
	if err != nil {
		Logger.Printf("Failed to send webhook err: %v\n", err)
	}
}
