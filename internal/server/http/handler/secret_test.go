package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

func TestAuth_Upload(t *testing.T) {
	secret := dto.SecretRequest{}
	reqBody, err := json.Marshal(secret)
	require.Nil(t, err, "Secret json encoding")

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		body  []byte
		setup func(t *testing.T) SecretService
		want  want
	}{
		{
			name: "success",
			body: reqBody,
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Save(gomock.All(), &secret).
					Return(nil)
				return service
			},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name: "ivalid_request_format",
			body: []byte("invalid json"),
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid request format",
			},
		},
		{
			name: "ivalid_secret_data",
			body: reqBody,
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Save(gomock.All(), &secret).
					Return(errors.ErrSecretInvalidData)
				return service
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid secret data",
			},
		},
		{
			name: "server_error",
			body: reqBody,
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Save(gomock.All(), &secret).
					Return(errors.ErrUnexpected)
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
			handler := NewSecret(service, logger)

			r := httptest.NewRequest(http.MethodPost, "/secret", bytes.NewBuffer(test.body))
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Upload(w, r)
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

func TestAuth_Get(t *testing.T) {
	secret := dto.SecretResponse{}
	respBody, err := json.Marshal(secret)
	require.Nil(t, err, "Secret json encoding")

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name     string
		secretID string
		setup    func(t *testing.T) SecretService
		want     want
	}{
		{
			name:     "succes",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Secret(gomock.All(), uint64(13)).
					Return(secret, nil)
				return service
			},
			want: want{
				code: http.StatusOK,
				body: string(respBody),
			},
		},
		{
			name:     "sbad_secret_id",
			secretID: "bad_id",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			want: want{
				code: http.StatusBadRequest,
				body: "invalid secret id",
			},
		},
		{
			name:     "secret_not_found",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Secret(gomock.All(), uint64(13)).
					Return(secret, errors.ErrSecretNotFound)
				return service
			},
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			name:     "server_error",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Secret(gomock.All(), uint64(13)).
					Return(secret, errors.ErrUnexpected)
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
			handler := NewSecret(service, logger)

			r := httptest.NewRequest(http.MethodGet, "/secret/"+test.secretID, nil)
			r.SetPathValue("id", test.secretID)

			w := httptest.NewRecorder()
			handler.Get(w, r)
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

func TestAuth_List(t *testing.T) {
	list := []dto.SecretInfo{{}}
	respBody, err := json.Marshal(list)
	require.Nil(t, err, "secret info list json encoding")

	type want struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		setup func(t *testing.T) SecretService
		want  want
	}{
		{
			name: "succes",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					InfoList(gomock.All()).
					Return(list, nil)
				return service
			},
			want: want{
				code: http.StatusOK,
				body: string(respBody),
			},
		},
		{
			name: "succes",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					InfoList(gomock.All()).
					Return(nil, errors.ErrUnexpected)
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
			handler := NewSecret(service, logger)

			r := httptest.NewRequest(http.MethodGet, "/secret", nil)

			w := httptest.NewRecorder()
			handler.List(w, r)
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
