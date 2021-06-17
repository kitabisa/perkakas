package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var unaryInfo = &grpc.UnaryServerInfo{
	FullMethod: "TestUnaryInterceptor",
}

func TestUnaryInterceptorWithoutRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.NotContains(t, out.String(), ctxkeys.CtxXKtbsRequestID)

		return nil, nil
	}

	ctx := context.Background()
	UnaryServerInterceptor(ctx, nil, unaryInfo, test)
}

func TestUnaryInterceptorWithRequestID(t *testing.T) {
	var out bytes.Buffer
	log.Logger = zerolog.New(&out).With().Caller().Logger()

	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {

		l := ctx.Value(ctxkeys.CtxLogger).(zerolog.Logger)

		l.Info().Err(errors.New("any-error")).Send()

		assert.Contains(t, out.String(), reqID)

		return nil, nil
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)
	UnaryServerInterceptor(ctx, nil, unaryInfo, test)
}
