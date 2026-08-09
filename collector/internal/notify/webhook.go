package notify

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pulseguard/shared/schema"
)

type Notifier struct {
	webhookURL string
	client     *http.Client
}

func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (n *Notifier) Enabled() bool {
	return n.webhookURL != ""
}

func (n *Notifier) SendEvent(e schema.Event) {
	if !n.Enabled() {
		return
	}

	payload, err := json.Marshal(map[string]string{
		"host_id":    e.HostID,
		"check_type": e.CheckType,
		"level":      string(e.Level),
		"message":    e.Message,
		"timestamp":  e.Timestamp,
	})
	if err != nil {
		log.Printf("webhook payload olusturulamadi: %v", err)
		return
	}

	go func() {
		resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(payload))
		if err != nil {
			log.Printf("webhook gonderilemedi: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("webhook gonderildi: %s (durum: %d)", e.EventID, resp.StatusCode)
	}()
}
