package interceptor

import (
	"context"
	"remote-force/server/entity"
	"remote-force/server/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserProvider interface {
	GetByID(id string) (entity.User, error)
}

func UserInfo(up UserProvider, s store.User) grpc.UnaryServerInterceptor {
	return excludePublicMethods(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			if userID, ok := ContextUserID(ctx); ok {
				user, err := up.GetByID(userID)
				if err != nil {
					return nil, status.Errorf(codes.Unauthenticated, "Retrieving user information is failed")
				}
				ctx = context.WithValue(ctx, contextKeyUser, user)
				return handler(ctx, req)
			}
			return nil, status.Errorf(codes.FailedPrecondition, "invalid context no user logged")
		})
}
