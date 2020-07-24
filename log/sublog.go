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
	return log.With().
		Str(internal.CtxXKtbsRequestID.String(), ctx.Value(internal.CtxXKtbsRequestID).(string)).
		Str("label", ctxName).
		Logger()
}
