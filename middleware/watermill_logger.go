package middleware

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog/log"
)

func WatermillLoggerMiddleware(h message.HandlerFunc) message.HandlerFunc {
	return func(message *message.Message) ([]*message.Message, error) {
		if message.Context().Value(ctxkeys.CtxLogger) != nil {
			return h(message)
		}
		logger := log.Logger
		processID := WatermillGetProcessID(message)
		if processID != "" {
			logger.With().Str(ctxkeys.CtxXKtbsRequestID.String(), processID)
		}
		ctx := context.WithValue(message.Context(), ctxkeys.CtxLogger, logger)
		message.SetContext(ctx)

		return h(message)
	}
}
