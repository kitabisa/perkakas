package middleware

import (
	"bytes"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatermillLoggerMiddleware(t *testing.T) {
	handler := WatermillLoggerMiddleware(func(msg *message.Message) ([]*message.Message, error) {
		return message.Messages{msg}, nil
	})

	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	msg := message.NewMessage("1", nil)

	producedMsgs, _ := handler(msg)
	processID := WatermillGetProcessID(producedMsgs[0])
	assert.Contains(t, out.String(), processID)
}
