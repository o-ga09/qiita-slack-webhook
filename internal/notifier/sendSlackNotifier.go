package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	ErrHTTPStatusNotOK = fmt.Errorf("HTTP status code is not OK")
)

// SlackMessage はSlackに送信するメッセージの構造体です
type SlackMessage struct {
	Text string `json:"text"`
}

func SendSlackNotification(ctx context.Context, message SlackMessage) error {
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if slackWebhookURL == "" {
		return fmt.Errorf("SLACK_WEBHOOK_URL is not set")
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling the message: %w", err)
	}

	resp, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error sending the message to Slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrHTTPStatusNotOK
	}
	return nil
}
