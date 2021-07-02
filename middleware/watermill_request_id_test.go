package middleware

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatermillProcessIDMiddleware(t *testing.T) {
	handler := WatermillProcessIDMiddleware(func(msg *message.Message) ([]*message.Message, error) {
		return message.Messages{msg}, nil
	})

	msg := message.NewMessage("1", nil)

	producedMsgs, _ := handler(msg)
	processID := WatermillGetProcessID(producedMsgs[0])
	fmt.Println("ProcessID", processID)
	assert.NotEmpty(t, processID)
}