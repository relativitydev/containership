package slack

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Client struct {
	WebhookURL string
}

func NewClient(webhookURL string) *Client {
	return &Client{
		WebhookURL: webhookURL,
	}
}

func (c *Client) SendMessage(message string) error {
	attachment := slack.Attachment{
		Text: message,
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook(c.WebhookURL, &msg)

	return errors.Wrap(err, "failed to send slack webhook")
}
