package unaryinterceptor

import (
	"context"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func LoggerToContextInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
	logger := log.Logger

	if reqID != "" {
		logger = log.With().Str(ctxkeys.CtxXKtbsRequestID.String(), reqID).Logger()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxLogger, logger)

	resp, err = handler(ctx, req)

	return
}
