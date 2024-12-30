package mw

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type key struct{}

func GetLogin(ctx context.Context) string {
	login, _ := ctx.Value(key{}).(string)
	return login
}

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}
	values := md.Get("x-auth")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no x-auth header")
	}

	ctx = context.WithValue(ctx, key{}, values[0])

	resp, err = handler(ctx, req)

	return resp, err
}
