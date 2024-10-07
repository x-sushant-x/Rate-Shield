package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackMessage represents the structure of the message payload
type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// SlackService is a struct that holds the Slack token and other configurations
type SlackService struct {
	Token   string
	Channel string
}

// NewSlackService is a constructor function to create a new SlackService instance
func NewSlackService(token, channel string) *SlackService {
	return &SlackService{
		Token:   token,
		Channel: channel,
	}
}

// SendSlackMessage sends a message to the configured Slack channel
func (s *SlackService) SendSlackMessage(message string) error {
	// Create the message payload
	slackMessage := SlackMessage{
		Channel: s.Channel,
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
	req.Header.Set("Authorization", "Bearer "+s.Token)

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
