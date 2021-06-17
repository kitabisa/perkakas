package unaryinterceptor

import (
	"context"

	"github.com/kitabisa/perkakas/v2/ctxkeys"
	"github.com/kitabisa/perkakas/v2/grpc/util"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GrpcRequestIDKey = "x-ktbs-request-id"
)

func ReqIDToContextInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqID, err := util.GetMetadata(ctx, GrpcRequestIDKey)
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
