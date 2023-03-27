package interceptor

import (
	"context"
	"remote-force/server/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Authentication(m *jwt.Manager) grpc.UnaryServerInterceptor {
	return excludePublicMethods(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			newCTX, err := authorize(ctx, m)
			if err != nil {
				return nil, err
			}
			return handler(newCTX, req)
		})
}

func authorize(ctx context.Context, m *jwt.Manager) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token := authHeader[0]

	c, err := m.Validate(token)

	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, err.Error())
	}
	ctx = context.WithValue(ctx, contextKeyUserID, c.ID)
	return ctx, nil
}
