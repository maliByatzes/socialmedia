package socialmedia

import "context"

type contextKey int

const userContextKey = contextKey(iota + 1)

func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}

func UserIDFromContext(ctx context.Context) uint {
	if user := UserFromContext(ctx); user != nil {
		return user.ID
	}
	return 0
}
