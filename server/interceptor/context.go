package interceptor

import (
	"context"
	"remote-force/server/entity"
)

type contextKey string

func (c contextKey) String() string {
	return "interceptor" + string(c)
}

var (
	contextKeyUserID = contextKey("user-id")
	contextKeyUser   = contextKey("user")
)

func ContextUserID(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(contextKeyUserID).(string)
	return tokenStr, ok
}

func ContextUser(ctx context.Context) (entity.User, bool) {
	usr, ok := ctx.Value(contextKeyUser).(entity.User)
	return usr, ok
}
