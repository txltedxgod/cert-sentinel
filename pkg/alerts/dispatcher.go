package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type AlertPayload struct {
	Target        string  `json:"target"`
	CommonName    string  `json:"common_name"`
	DaysRemaining float64 `json:"days_remaining"`
	Timestamp     string  `json:"timestamp"`
	Message       string  `json:"message"`
}

type Dispatcher struct {
	webhookURL string
	client     *http.Client
}

func NewDispatcher(webhookURL string) *Dispatcher {
	return &Dispatcher{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (d *Dispatcher) SendAlert(target, commonName string, daysRemaining float64) {
	if d.webhookURL == "" {
		return
	}

	payload := AlertPayload{
		Target:        target,
		CommonName:    commonName,
		DaysRemaining: daysRemaining,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Message:       fmt.Sprintf("⚠️ SSL/TLS Certificate for %s expires in %.1f days!", target, daysRemaining),
	}

	body, _ := json.Marshal(payload)
	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[Alert Error] Failed to send webhook alert for %s: %v\n", target, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("[Alert Sent] Successfully dispatched expiry warning for %s\n", target)
}
