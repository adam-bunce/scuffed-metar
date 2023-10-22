package globals

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var ServerPort = 80
var webhook = os.Getenv("WEBHOOK_URL")

func SendWebhook(message string) {
	_, err := http.Post(webhook, "application/json", strings.NewReader(fmt.Sprintf(`{"content": "%s"}`, message)))
	if err != nil {
		log.Printf("Failed to send webhook err: %v\n", err)
	}
}
