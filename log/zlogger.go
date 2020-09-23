package log

import (
	"context"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog"
)

// Zlogger get zerolog sublogger from context
func Zlogger(ctx context.Context) *zerolog.Logger {
	var logger *zerolog.Logger
	if ctx.Value(ctxkeys.CtxLogger) != nil {
		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)
		logger = &l
	}

	return logger
}
