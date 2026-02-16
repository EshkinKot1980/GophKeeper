package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

func TestAuth_Register(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	salt := "duBXKxwaWXfhgXBQrrwdtQ"
	successBody := `{"token":"` + token + `","encr_salt":"` + salt + `"}`

	errLoginTooLong := fmt.Errorf(
		"%w: login too long, max %d characters",
		errors.ErrAuthInvalidCredentials,
		entity.UserMaxLoginLen,
	)

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		body  string
		setup func(t *testing.T) AuthService
		want  want
	}{
		{
			name: "success",
			body: `{"login":"testLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Login: "testLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: salt}, nil)
				return service
			},
			want: want{
				code: http.StatusOK,
				body: successBody,
			},
		},
		{
			name: "negative_bad_json",
			body: `not valid jsson`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockAuthService(ctrl)
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid credentials format",
			},
		},
		{
			name: "negative_without_login",
			body: `{"password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Password: "t1estP5assword"}).
					Return(dto.AuthResponse{}, errors.ErrAuthInvalidCredentials)
				return service
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid credentials",
			},
		},
		{
			name: "negative_without_password",
			body: `{"login":"testLogin"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Login: "testLogin"}).
					Return(dto.AuthResponse{}, errors.ErrAuthInvalidCredentials)
				return service
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid credentials",
			},
		},
		{
			name: "negative_login_too_long",
			body: `{"login":"veryLongLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Login: "veryLongLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{}, errLoginTooLong)
				return service
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid credentials: login too long, max 64 characters",
			},
		},
		{
			name: "negative_user_already_exists",
			body: `{"login":"testLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Login: "testLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{}, errors.ErrAuthUserAlreadyExists)
				return service
			},
			want: want{
				code: http.StatusConflict,
				body: "user already exists",
			},
		},
		{
			name: "negative_server_error",
			body: `{"login":"testLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All(), dto.Credentials{Login: "testLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{}, errors.ErrUnexpected)
				return service
			},
			want: want{
				code: http.StatusInternalServerError,
				body: statusText500,
			},
		},
	}

	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := test.setup(t)
			handler := NewAuth(service, logger)

			reqBody := []byte(test.body)
			r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Register(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode, "Response status code")

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := strings.TrimSuffix(string(resBody), "\n")
			assert.Equal(t, test.want.body, body, "Response body")
		})
	}
}

func TestAuth_Login(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	salt := "duBXKxwaWXfhgXBQrrwdtQ"
	successBody := `{"token":"` + token + `","encr_salt":"` + salt + `"}`

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		body  string
		setup func(t *testing.T) AuthService
		want  want
	}{
		{
			name: "success",
			body: `{"login":"testLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Login(gomock.All(), dto.Credentials{Login: "testLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: salt}, nil)
				return service
			},
			want: want{
				code: http.StatusOK,
				body: successBody,
			},
		},
		{
			name: "negative_bad_json",
			body: `not valid jsson`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockAuthService(ctrl)
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid credentials format",
			},
		},
		{
			name: "negative_bad_login_or_password",
			body: `{"login":"badLogin", "password":"orPassword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Login(gomock.All(), dto.Credentials{Login: "badLogin", Password: "orPassword"}).
					Return(dto.AuthResponse{}, errors.ErrAuthInvalidCredentials)
				return service
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "",
			},
		},
		{
			name: "negative_server_error",
			body: `{"login":"testLogin", "password":"t1estP5assword"}`,
			setup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Login(gomock.All(), dto.Credentials{Login: "testLogin", Password: "t1estP5assword"}).
					Return(dto.AuthResponse{}, errors.ErrUnexpected)
				return service
			},
			want: want{
				code: http.StatusInternalServerError,
				body: statusText500,
			},
		},
	}

	ctrl := gomock.NewController(t)
	logger := mocks.NewMockLogger(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := test.setup(t)
			handler := NewAuth(service, logger)

			reqBody := []byte(test.body)
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Login(w, r)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode, "Response status code")

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			body := strings.TrimSuffix(string(resBody), "\n")
			assert.Equal(t, test.want.body, body, "Response body")
		})
	}
}
