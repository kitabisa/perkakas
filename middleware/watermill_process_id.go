package middleware

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/kitabisa/perkakas/v2/ctxkeys"
)

func WatermillGetProcessID(message *message.Message) string {
	ctx := message.Context()
	if processID, exist := ctx.Value(ctxkeys.CtxWatermillProcessID).(string); exist {
		return processID
	}
	return ""
}

func WatermillSetProcessID(message *message.Message) {
	processID := uuid.New().String()
	ctx := context.WithValue(message.Context(), ctxkeys.CtxWatermillProcessID, processID)
	message.SetContext(ctx)
}

func WatermillProcessIDMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(message *message.Message) ([]*message.Message, error) {
		processID := WatermillGetProcessID(message)
		if processID != "" {
			return h(message)
		}

		WatermillSetProcessID(message)

		return h(message)
	}
}
