package models

type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
