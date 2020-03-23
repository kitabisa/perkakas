package slack

import (
	"errors"
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
)

// WebHook defines simple webhook instance
type WebHook struct {
	URL              string
	IsChannelMention bool
	Attachment       []slack.Attachment
	Text             string
}

// NewWebhook will create simple webhook instance
func NewWebhook(url string) *WebHook {
	defaultColor := "#FF5733"
	attachment1 := slack.Attachment{}
	attachment1.Color = &defaultColor

	webhook := &WebHook{
		URL:        url,
		Attachment: []slack.Attachment{attachment1},
	}

	return webhook
}

// AddField will new field to slack attachment
func (w *WebHook) AddField(title, value string) {
	w.Attachment[0].AddField(slack.Field{Title: title, Value: value})
}

// AddText will add text to slack message
func (w *WebHook) AddText(message string) {
	w.Text = message
}

// SetChannelMention will set whether @channel is mentioned or not
func (w *WebHook) SetChannelMention(flag bool) {
	w.IsChannelMention = flag
}

// SetColor will set attachment color in #FFFFFF hex representation
func (w *WebHook) SetColor(color string) {
	w.Attachment[0].Color = &color
}

// Send will send slack notification using webhook
func (w *WebHook) Send() error {
	if w.Text == "" {
		return errors.New("error: slack message cannot be empty")
	}

	textMsg := w.Text
	if w.IsChannelMention {
		textMsg = fmt.Sprintf("<!channel> %s", w.Text)
	}

	payload := slack.Payload{
		Text:        textMsg,
		Attachments: w.Attachment,
	}

	err := slack.Send(w.URL, "", payload)
	if len(err) > 0 {
		return fmt.Errorf("error: %s", err)
	}
	return nil
}
