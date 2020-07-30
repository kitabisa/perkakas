package log

import (
	"context"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetSublogger get zerolog sublogger
// WIP: middleware for set X-Ktbs-Request-ID to context
func GetSublogger(ctx context.Context, ctxName string) zerolog.Logger {
	reqID := ""
	if ctx.Value(ctxkeys.CtxXKtbsRequestID) != nil {
		reqID = ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
	}

	return log.With().
		Str(ctxkeys.CtxXKtbsRequestID.String(), reqID).
		Str("label", ctxName).
		Logger()
}
