package requestid

import (
	"context"
	"testing"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var unaryInfo = &grpc.UnaryServerInfo{
	FullMethod: "TestUnaryInterceptor",
}

func TestUnaryInterceptorWithoutRequestID(t *testing.T) {

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID).(string)
		assert.NotEmpty(t, requestID)

		return requestID, nil
	}

	UnaryServerInterceptor(context.Background(), nil, unaryInfo, test)

}

func TestUnaryInterceptorWithRequestID(t *testing.T) {
	reqID := uuid.NewV4().String()

	test := func(ctx context.Context, req interface{}) (interface{}, error) {
		requestID := ctx.Value(ctxkeys.CtxXKtbsRequestID)

		assert.Equal(t, reqID, requestID)

		return requestID, nil
	}

	ctx := context.Background()
	md := metadata.Pairs(GrpcRequestIDKey, reqID)
	ctx = metadata.NewIncomingContext(ctx, md)
	UnaryServerInterceptor(ctx, nil, unaryInfo, test)
}
