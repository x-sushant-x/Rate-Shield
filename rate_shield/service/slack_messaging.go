package service

import (
	"bytes"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/utils"
)

const (
	SLACK_SEND_MESSAGE_ENDPOINT = "https://slack.com/api/chat.postMessage"
)

type SlackService struct {
	Token   string
	Channel string
}

func NewSlackService(token, channel string) *SlackService {
	return &SlackService{
		Token:   token,
		Channel: channel,
	}
}

func (s *SlackService) SendSlackMessage(msg string) error {
	message := buildSlackMessageObject(s.Channel, msg)

	messageBytes, err := utils.MarshalJSON(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, SLACK_SEND_MESSAGE_ENDPOINT, bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Err(err).Msgf("error creating request: %s", err)
		return err
	}
	s.setRequestHeaders(req)

	return s.sendRequestToSlackAPI(req)
}

func buildSlackMessageObject(channel, msg string) models.SlackMessage {
	message := models.SlackMessage{
		Channel: channel,
		Text:    msg,
	}

	return message
}

func (s *SlackService) setRequestHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token)
}

func (s *SlackService) sendRequestToSlackAPI(req *http.Request) error {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msgf("error sending request: %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Err(err).Msgf("received non-OK response from Slack: %s", resp.Status)
		return err
	}

	return nil
}
