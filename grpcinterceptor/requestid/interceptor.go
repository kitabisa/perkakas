package requestid

import (
	"context"
	"errors"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	GrpcRequestIDKey = "x-ktbs-request-id"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqID, err := getRequestID(ctx, GrpcRequestIDKey)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if reqID == "" {
		reqID = uuid.NewV4().String()
	}

	ctx = context.WithValue(ctx, ctxkeys.CtxXKtbsRequestID, reqID)

	resp, err = handler(ctx, req)

	return
}

func getRequestID(ctx context.Context, key string) (val string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = errors.New("failed to retrieve metadata")
	}

	v := md[key]

	if len(v) == 0 {
		return
	}

	val = v[0]

	return
}
