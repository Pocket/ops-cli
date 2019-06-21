package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Attachment struct {
	Color     string `json:"text"`
	Fallback  string `json:"text"`
	Title     string `json:"text"`
	TitleLink string `json:"text"`
}

type Request struct {
	Text        string `json:"text"`
	Username    string `json:"text"`
	Channel     string `json:"text"`
	IconEmoji   string `json:"text"`
	Attachments []Attachment
}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func (requestBody *Request) SendSlackNotification(webhookUrl string) error {

	slackBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

func NewSlackRequest(username string, channel string, iconEmoji string, text string, color string, fallback string, title string, titleLink string) *Request {
	return &Request{
		Username:  username,
		Channel:   channel,
		IconEmoji: iconEmoji,
		Text:      text,
		Attachments: []Attachment{
			{
				Color:     color,
				Fallback:  fallback,
				Title:     title,
				TitleLink: titleLink,
			},
		},
	}
}
