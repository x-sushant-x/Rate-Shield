package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// SlackMessage represents the structure of the message payload
type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// SendSlackMessage sends a message to a Slack channel
func SendSlackMessage(message, channel string) error {
	// Get the Slack OAuth token from environment variable
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("slack token is not set")
	}

	// Create the message payload
	slackMessage := SlackMessage{
		Channel: channel,
		Text:    message,
	}

	// Marshal the message to JSON
	messageBytes, err := json.Marshal(slackMessage)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	// Create a POST request to the Slack API
	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(messageBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+slackToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response from Slack
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Slack: %s", resp.Status)
	}

	return nil
}
