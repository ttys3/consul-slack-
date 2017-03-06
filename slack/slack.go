package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
)

// New creates new slack client
func New(cfg *Config) (*Slack, error) {
	if cfg == nil {
		panic("cfg is nil")
	}

	return &Slack{cfg: cfg}, nil
}

// Config is slack config
type Config struct {
	// WebhookURL is webhook url
	WebhookURL string

	// Chanel is slack channel name, starts with #
	Channel string

	// Username of user that messages are sent on behalf of
	Username string

	// IconURL is avatar url
	IconURL string
}

// Slack is slack client
type Slack struct {
	cfg *Config
}

// payload is data that is sent to webhook url
type payload struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	IconURL     string       `json:"icon_url"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Color string `json:"color"`
	Text  string `json:"text"`
}

// Danger sends a danger message
func (s *Slack) Danger(msg string, v ...interface{}) error {
	return s.send("danger", msg, v...)
}

// Good sends a good message
func (s *Slack) Good(msg string, v ...interface{}) error {
	return s.send("good", msg, v...)
}

// send sends message to webhook url
func (s *Slack) send(color, msg string, v ...interface{}) error {
	if s.cfg.WebhookURL == "" {
		return nil
	}

	b, err := json.Marshal(&payload{
		Channel:  s.cfg.Channel,
		Username: s.cfg.Username,
		IconURL:  s.cfg.IconURL,
		Attachments: []attachment{
			{
				Color: color,
				Text:  fmt.Sprintf(msg, v...),
			},
		},
	})

	if err != nil {
		return err
	}

	_, err = http.Post(s.cfg.WebhookURL, "application/json", bytes.NewReader(b))
	return err
}
