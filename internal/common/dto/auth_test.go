package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentials_Validate(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  string
	}{
		{
			name:     "succes",
			login:    "test13",
			password: "password13",
		},
		{
			name:     "invalid_login",
			login:    "t1",
			password: "password13",
			wantErr:  "login too short (min 3 chars)",
		},
		{
			name:     "invalid_password",
			login:    "test13",
			password: "p1",
			wantErr:  "password too short (min 8 chars)",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotErr string
			cr := Credentials{Login: test.login, Password: test.password}

			err := cr.Validate()
			if err != nil {
				gotErr = err.Error()
			}

			assert.Equal(t, test.wantErr, gotErr, "Validation error")
		})
	}
}

func TestCredentials_ValidateLogin(t *testing.T) {
	tooLongLogin := "l"
	for range CredentialsLoginMaxLen {
		tooLongLogin += "l"
	}

	tests := []struct {
		name    string
		login   string
		wantErr string
	}{
		{
			name:  "succes",
			login: "test13",
		},
		{
			name:    "too_short_login",
			login:   "t1",
			wantErr: "login too short (min 3 chars)",
		},
		{
			name:    "too_long_login",
			login:   tooLongLogin,
			wantErr: "login too long (max 64 chars)",
		},
		{
			name:    "invalid_characters",
			login:   "#$%",
			wantErr: "login must contain letters or digits",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotErr string
			cr := Credentials{Login: test.login}

			err := cr.ValidateLogin()
			if err != nil {
				gotErr = err.Error()
			}

			assert.Equal(t, test.wantErr, gotErr, "Validation error")
		})
	}
}

func TestCredentials_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  string
	}{
		{
			name:     "succes",
			password: "password13",
		},
		{
			name:     "too_short_password",
			password: "pass123",
			wantErr:  "password too short (min 8 chars)",
		},
		{
			name:     "without_numbers",
			password: "password",
			wantErr:  "password must contain both letters and digits",
		},
		{
			name:     "without_numbers",
			password: "12345678",
			wantErr:  "password must contain both letters and digits",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotErr string
			cr := Credentials{Password: test.password}

			err := cr.ValidatePassword()
			if err != nil {
				gotErr = err.Error()
			}

			assert.Equal(t, test.wantErr, gotErr, "Validation error")

		})
	}
}
