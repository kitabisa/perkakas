package log

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

var (
	CtxXKtbsRequestID contextKey = "X-Ktbs-Request-ID"
)

func (c contextKey) String() string {
	return string(c)
}

// GetSublogger get zerolog sublogger
func GetSublogger(ctx context.Context, ctxName string) zerolog.Logger {
	return log.With().
		Str(CtxXKtbsRequestID.String(), ctx.Value(CtxXKtbsRequestID).(string)).
		Str("label", ctxName).
		Logger()
}
