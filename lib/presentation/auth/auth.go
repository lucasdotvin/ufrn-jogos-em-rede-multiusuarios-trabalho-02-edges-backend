package auth

import "context"

const (
	userUuidKey = "user_uuid"
)

func WithUserUuid(ctx context.Context, userUuid string) context.Context {
	return context.WithValue(ctx, userUuidKey, userUuid)
}

func GetUserUuid(ctx context.Context) (string, bool) {
	userUuid, ok := ctx.Value(userUuidKey).(string)

	return userUuid, ok
}
