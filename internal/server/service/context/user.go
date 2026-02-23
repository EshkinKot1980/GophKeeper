// Пвкет context содержит гетеры сетеры для контекста сервисов
package context

import (
	"context"
	"fmt"
)

type сontextKey string

const keyUserID сontextKey = "userID"

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

func UserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(keyUserID).(string)
	if !ok {
		return "", fmt.Errorf("failed to get string value for %s key", keyUserID)
	}
	return userID, nil
}
