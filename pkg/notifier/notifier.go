package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Notifier handles sending alerts to external services (Discord/Slack).
type Notifier struct {
	webhookURL string
	client     *http.Client
	// rateLimitCache stores the last alert time for an IP to prevent spam.
	rateLimitCache    sync.Map
	rateLimitDuration time.Duration
}

// New creates a new Notifier instance.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL:        webhookURL,
		client:            &http.Client{Timeout: 5 * time.Second},
		rateLimitDuration: 1 * time.Hour,
	}
}

// SendAlert sends a notification to the configured webhook.
// It includes rate limiting logic to avoid spamming for the same IP.
func (n *Notifier) SendAlert(trapName, remoteIP, userAgent string) {
	if n.webhookURL == "" {
		return
	}

	// Rate Limiting: Check if we've seen this IP recently
	if lastSeen, ok := n.rateLimitCache.Load(remoteIP); ok {
		if time.Since(lastSeen.(time.Time)) < n.rateLimitDuration {
			return // Skip alert
		}
	}
	// Update last seen time
	n.rateLimitCache.Store(remoteIP, time.Now())

	// Run async to not block the trap handler
	go n.send(trapName, remoteIP, userAgent)
}

func (n *Notifier) send(trapName, remoteIP, userAgent string) {
	// Construct payload compatible with Discord (content) and Slack (text)
	// We'll use a generic map to send both fields if needed, or just "content" for Discord.
	// For this implementation, we target Discord's format.
	msg := fmt.Sprintf("🚨 **Trap Triggered!**\n**Trap:** `%s`\n**IP:** `%s`\n**UA:** `%s`", trapName, remoteIP, userAgent)

	payload := map[string]string{
		"content": msg, // Discord
		"text":    msg, // Slack
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal alert payload")
		return
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send webhook alert")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Error().Int("status", resp.StatusCode).Msg("Webhook returned error status")
	} else {
		log.Debug().Str("trap", trapName).Msg("Alert sent successfully")
	}
}
