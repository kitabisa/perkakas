package log

import (
	"context"

	"github.com/kitabisa/perkakas/v2/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetSublogger get zerolog sublogger
// WIP: middleware for set X-Ktbs-Request-ID to context
func GetSublogger(ctx context.Context, ctxName string) zerolog.Logger {
	reqID := ""
	if ctx.Value(internal.CtxXKtbsRequestID) != nil {
		reqID = ctx.Value(internal.CtxXKtbsRequestID).(string)
	}

	return log.With().
		Str(internal.CtxXKtbsRequestID.String(), reqID).
		Str("label", ctxName).
		Logger()
}
