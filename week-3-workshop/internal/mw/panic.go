package mw

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func Panic(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()
	return handler(ctx, req)
}
