package util

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

func GetMetadata(ctx context.Context, key string) (val string, err error) {
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
