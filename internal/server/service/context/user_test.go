package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUserID(t *testing.T) {
	userID := "d7d81ca8-8b0b-496e-abbd-fd522245c975"

	ctx := SetUserID(context.Background(), userID)
	gotUserID, ok := ctx.Value(keyUserID).(string)

	assert.True(t, ok, "Value string assertion")
	assert.Equal(t, userID, gotUserID, "UserID value comparation")
}

func TestUserID(t *testing.T) {
	userID := "d7d81ca8-8b0b-496e-abbd-fd522245c975"

	type want struct {
		userID string
		hasErr bool
	}

	tests := []struct {
		name string
		ctx  context.Context
		want want
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), keyUserID, userID),
			want: want{userID: userID},
		},
		{
			name: "without_key",
			ctx:  context.Background(),
			want: want{hasErr: true},
		},
		{
			name: "bad_type",
			ctx:  context.WithValue(context.Background(), keyUserID, 13),
			want: want{hasErr: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUserID, err := UserID(test.ctx)
			gorErr := err != nil
			assert.Equal(t, test.want.userID, gotUserID, "UserID value comparation")
			assert.Equal(t, test.want.hasErr, gorErr, "Error presence comparation")
		})
	}
}
