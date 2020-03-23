package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/stretchr/testify/assert"
)

func emptyTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
}

func TestWebhookCase1(t *testing.T) {
	slackMsg := NewWebhook("https://my-webhook-url")
	slackMsg.AddText("my slack text")
	slackMsg.AddField("name", "mas berto")
	slackMsg.AddField("phonenumber", "081288889999")
	slackMsg.SetChannelMention(true)

	attachment1 := slack.Attachment{}
	attachment1.AddField(slack.Field{Title: "name", Value: "mas berto"})
	attachment1.AddField(slack.Field{Title: "phonenumber", Value: "081288889999"})

	color := "#FF5733"
	attachment1.Color = &color

	expected := WebHook{
		URL:              "https://my-webhook-url",
		Attachment:       []slack.Attachment{attachment1},
		Text:             "my slack text",
		IsChannelMention: true,
	}

	assert.Equal(t, expected, slackMsg)
}

func TestWebhookCase2(t *testing.T) {
	slackMsg := NewWebhook("https://my-webhook-url")
	slackMsg.AddText("my slack text")
	slackMsg.SetColor("#000000")

	color := "#000000"
	attachment1 := slack.Attachment{}
	attachment1.Color = &color

	expected := WebHook{
		URL:              "https://my-webhook-url",
		Attachment:       []slack.Attachment{attachment1},
		Text:             "my slack text",
		IsChannelMention: false,
	}

	assert.Equal(t, expected, slackMsg)
}

func TestWebhookSend_TextEmpty(t *testing.T) {
	emptyServer := emptyTestServer()

	slackMsg := NewWebhook(emptyServer.URL)
	slackMsg.AddField("name", "mas berto")
	slackMsg.AddField("phonenumber", "081288889999")
	slackMsg.SetChannelMention(true)

	err := slackMsg.Send()
	assert.NotNil(t, err)
}

func TestWebhookSend(t *testing.T) {
	emptyServer := emptyTestServer()

	slackMsg := NewWebhook(emptyServer.URL)
	slackMsg.AddText("my slack text")
	slackMsg.AddField("name", "mas berto")
	slackMsg.AddField("phonenumber", "081288889999")
	slackMsg.SetChannelMention(true)

	err := slackMsg.Send()
	assert.Equal(t, nil, err)
}
