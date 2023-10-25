package globals

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var ServerPort = 8080
var webhook = os.Getenv("WEBHOOK_URL")
var Logger = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

func SendWebhook(message string) {
	_, err := http.Post(webhook, "application/json", strings.NewReader(fmt.Sprintf(`{"content": "%s"}`, message)))
	if err != nil {
		Logger.Printf("Failed to send webhook err: %v\n", err)
	}
}
