package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

const (
	login_method = "/V1.Remote/Login"
)

func excludeMethods(methods []string) func(grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			for _, method := range methods {
				if info.FullMethod == method {
					return handler(ctx, req)
				}
			}
			return next(ctx, req, info, handler)
		}
	}
}

var excludeLogin = excludeMethods([]string{login_method})
