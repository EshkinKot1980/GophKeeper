package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Register(t *testing.T) {
	credentials := dto.Credentials{Login: "test", Password: "password13"}
	reqBody, err := json.Marshal(credentials)
	require.Nil(t, err, "Credentials json encoding")

	authResp := dto.AuthResponse{Token: "token", EncrSalt: "encryption salt"}
	respBody, err := json.Marshal(authResp)
	require.Nil(t, err, "Auth response json encoding")

	type want struct {
		authResp dto.AuthResponse
		err      error
	}

	tests := []struct {
		name     string
		netError bool
		respCode int
		respBody []byte
		want     want
	}{
		{
			name:     "succes",
			respCode: http.StatusOK,
			respBody: respBody,
			want: want{
				authResp: authResp,
			},
		},
		{
			name:     "network_error",
			netError: true,
			want: want{
				err: ErrRegistrationFailed,
			},
		},
		{
			name:     "registration_failed",
			respCode: http.StatusBadRequest,
			want: want{
				err: ErrRegistrationFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, RegisterPath, r.RequestURI, "Request URI")
				assert.Equal(t, http.MethodPost, r.Method, "Request Method")

				body, err := io.ReadAll(r.Body)
				require.Nil(t, err, "Read response body")
				assert.Equal(t, reqBody, body, "Request body")

				if test.respCode != http.StatusOK {
					w.WriteHeader(test.respCode)
					return
				}

				w.Header().Set("Content-Type", ContentType)
				w.WriteHeader(test.respCode)
				_, err = w.Write(test.respBody)
				require.Nil(t, err, "Write response body")
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			client := NewClient(server.URL, true)
			if test.netError {
				server.Close()
			}

			resp, err := client.Register(credentials)
			assert.ErrorIs(t, err, test.want.err, "Register error")
			if err != nil {
				return
			}
			assert.Equal(t, test.want.authResp, resp, "Auth response")

		})
	}

}

func TestClient_Login(t *testing.T) {
	credentials := dto.Credentials{Login: "test", Password: "password13"}
	reqBody, err := json.Marshal(credentials)
	require.Nil(t, err, "Credentials json encoding")

	authResp := dto.AuthResponse{Token: "token", EncrSalt: "encryption salt"}
	respBody, err := json.Marshal(authResp)
	require.Nil(t, err, "Auth response json encoding")

	type want struct {
		authResp dto.AuthResponse
		err      error
	}

	tests := []struct {
		name     string
		netError bool
		respCode int
		respBody []byte
		want     want
	}{
		{
			name:     "succes",
			respCode: http.StatusOK,
			respBody: respBody,
			want: want{
				authResp: authResp,
			},
		},
		{
			name:     "network_error",
			netError: true,
			want: want{
				err: ErrLoginFailed,
			},
		},
		{
			name:     "login_failed",
			respCode: http.StatusBadRequest,
			want: want{
				err: ErrLoginFailed,
			},
		},
		{
			name:     "server_error",
			respCode: http.StatusInternalServerError,
			want: want{
				err: ErrLoginFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, LoginPath, r.RequestURI, "Request URI")
				assert.Equal(t, http.MethodPost, r.Method, "Request Method")

				body, err := io.ReadAll(r.Body)
				require.Nil(t, err, "Read response body")
				assert.Equal(t, reqBody, body, "Request body")

				if test.respCode != http.StatusOK {
					w.WriteHeader(test.respCode)
					return
				}

				w.Header().Set("Content-Type", ContentType)
				w.WriteHeader(test.respCode)
				_, err = w.Write(test.respBody)
				require.Nil(t, err, "Write response body")
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			client := NewClient(server.URL, true)
			if test.netError {
				server.Close()
			}

			resp, err := client.Login(credentials)
			assert.ErrorIs(t, err, test.want.err, "Login error")
			if err != nil {
				return
			}
			assert.Equal(t, test.want.authResp, resp, "Auth response")

		})
	}
}

func TestClient_Upload(t *testing.T) {
	secretRequest := dto.SecretRequest{}
	reqBody, err := json.Marshal(secretRequest)
	require.Nil(t, err, "Secret request json encoding")

	secretResp := dto.SecretResponse{}
	respBody, err := json.Marshal(secretResp)
	require.Nil(t, err, "Auth response json encoding")

	tests := []struct {
		name     string
		netError bool
		respCode int
		repBody  []byte
		wantErr  error
	}{
		{
			name:     "succes",
			respCode: http.StatusOK,
			repBody:  respBody,
		},
		{
			name:     "network_error",
			netError: true,
			wantErr:  ErrSecretSendFailed,
		},
		{
			name:     "secret_upload_error",
			respCode: http.StatusBadRequest,
			wantErr:  ErrSecretSendFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, SecretPath, r.RequestURI, "Request URI")
				assert.Equal(t, http.MethodPost, r.Method, "Request Method")
				assert.Equal(t, "Bearer token", r.Header.Get("Authorization"), "Authorization header")

				body, err := io.ReadAll(r.Body)
				require.Nil(t, err, "Read response body")
				assert.Equal(t, reqBody, body, "Request body")

				w.WriteHeader(test.respCode)

				if test.respCode != http.StatusOK {
					_, err = w.Write([]byte("error description"))
					require.Nil(t, err, "Write response body")
				}
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			client := NewClient(server.URL, true)
			if test.netError {
				server.Close()
			}

			err := client.Upload(secretRequest, "token")
			assert.ErrorIs(t, err, test.wantErr, "Upload error")

		})
	}
}

func TestClient_Retrieve(t *testing.T) {
	secretResp := dto.SecretResponse{}
	respBody, err := json.Marshal(secretResp)
	require.Nil(t, err, "Secret response json encoding")

	type want struct {
		secretResp dto.SecretResponse
		err        error
	}

	tests := []struct {
		name     string
		netError bool
		respCode int
		respBody []byte
		want     want
	}{
		{
			name:     "succes",
			respCode: http.StatusOK,
			respBody: respBody,
			want: want{
				secretResp: secretResp,
			},
		},
		{
			name:     "network_error",
			netError: true,
			want: want{
				err: ErrSecretRetrieveFailed,
			},
		},
		{
			name:     "unauthorized",
			respCode: http.StatusUnauthorized,
			want: want{
				err: ErrSecretRetrieveFailed,
			},
		},
		{
			name:     "bad_request",
			respCode: http.StatusBadRequest,
			want: want{
				err: ErrSecretRetrieveFailed,
			},
		},
		{
			name:     "not_found",
			respCode: http.StatusNotFound,
			want: want{
				err: ErrSecretRetrieveFailed,
			},
		},
		{
			name:     "internal_server_error",
			respCode: http.StatusInternalServerError,
			want: want{
				err: ErrSecretRetrieveFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, SecretPath+"/13", r.RequestURI, "Request URI")
				assert.Equal(t, http.MethodGet, r.Method, "Request Method")
				assert.Equal(t, "Bearer token", r.Header.Get("Authorization"), "Authorization header")

				if test.respCode != http.StatusOK {
					w.WriteHeader(test.respCode)
					return
				}

				w.Header().Set("Content-Type", ContentType)
				w.WriteHeader(test.respCode)
				_, err = w.Write(respBody)
				require.Nil(t, err, "Write response body")
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			client := NewClient(server.URL, true)
			if test.netError {
				server.Close()
			}

			resp, err := client.Retrieve(13, "token")
			assert.ErrorIs(t, err, test.want.err, "Retrieve error")
			if err == nil {
				assert.Equal(t, test.want.secretResp, resp, "Secret response")
			}
		})
	}
}

func TestClient_InfoList(t *testing.T) {
	infoList := []dto.SecretInfo{{}}
	respBody, err := json.Marshal(infoList)
	require.Nil(t, err, "Auth response json encoding")

	type want struct {
		list []dto.SecretInfo
		err  error
	}

	tests := []struct {
		name     string
		netError bool
		respCode int
		respBody []byte
		want     want
	}{
		{
			name:     "succes",
			respCode: http.StatusOK,
			respBody: respBody,
			want: want{
				list: infoList,
			},
		},
		{
			name:     "network_error",
			netError: true,
			want: want{
				err: ErrSecretInfoListFailed,
			},
		},
		{
			name:     "unauthorized",
			respCode: http.StatusUnauthorized,
			want: want{
				err: ErrSecretInfoListFailed,
			},
		},
		{
			name:     "internal_server_error",
			respCode: http.StatusInternalServerError,
			want: want{
				err: ErrSecretInfoListFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, SecretPath, r.RequestURI, "Request URI")
				assert.Equal(t, http.MethodGet, r.Method, "Request Method")
				assert.Equal(t, "Bearer token", r.Header.Get("Authorization"), "Authorization header")

				if test.respCode != http.StatusOK {
					w.WriteHeader(test.respCode)
					return
				}

				w.Header().Set("Content-Type", ContentType)
				w.WriteHeader(test.respCode)
				_, err = w.Write(respBody)
				require.Nil(t, err, "Write response body")
			}

			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			client := NewClient(server.URL, true)
			if test.netError {
				server.Close()
			}

			list, err := client.InfoList("token")
			assert.ErrorIs(t, err, test.want.err, "Retrieve error")
			if err == nil {
				assert.Equal(t, test.want.list, list, "Secret info list")
			}
		})
	}
}
