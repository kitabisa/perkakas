package logger

import (
	"context"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	logger := log.Logger

	reqID, ok := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
	if ok {
		logger = log.With().Str(ctxkeys.CtxXKtbsRequestID.String(), reqID).Logger()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxLogger, logger)

	resp, err = handler(ctx, req)

	return
}
