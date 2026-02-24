package handler

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler/mocks"
)

func Test_jsonWriter_write(t *testing.T) {
	type want struct {
		code   int
		header string
		body   string
	}

	validLSON := `{"token":"TokenString","encr_salt":"EncryptionSaltString"}`

	tests := []struct {
		name       string
		witer      *testResponseWriter
		logger     func(t *testing.T) Logger
		value      any
		valueName  string
		stasusCode int
		want       want
	}{
		{
			name:  "positive",
			witer: newTestResponseWriter(false),
			logger: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				return mocks.NewMockLogger(ctrl)
			},
			value: dto.AuthResponse{
				Token:    "TokenString",
				EncrSalt: "EncryptionSaltString",
			},
			valueName:  "AuthResponse",
			stasusCode: http.StatusOK,
			want: want{
				code:   http.StatusOK,
				header: "application/json",
				body:   validLSON,
			},
		},
		{
			name:  "negative_failed_to_json_encode",
			witer: newTestResponseWriter(false),
			logger: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().Error("failed to encode order to json", gomock.All()).
					Times(1)
				return logger
			},
			value:      math.Inf(-1),
			valueName:  "order",
			stasusCode: http.StatusOK,
			want: want{
				code:   http.StatusInternalServerError,
				header: "text/plain",
				body:   statusText500,
			},
		},
		{
			name:  "negative_failed_to_json_encode",
			witer: newTestResponseWriter(true),
			logger: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().Error("failed to write body", gomock.All()).
					Times(1)
				return logger
			},
			value: dto.AuthResponse{
				Token:    "TokenString",
				EncrSalt: "EncryptionSaltString",
			},
			valueName:  "AuthResponse",
			stasusCode: http.StatusOK,
			want: want{
				code:   http.StatusOK,
				header: "application/json",
				body:   "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := test.witer
			l := test.logger(t)

			jw := newJSONwriter(w, l)
			jw.write(test.value, test.valueName, test.stasusCode)

			assert.Equal(t, test.want.code, w.Code, "Response status code")
			assert.Contains(t, w.Header().Get("Content-Type"), test.want.header, "Response content type")
			resBody := strings.TrimSuffix(string(w.Body.String()), "\n")
			assert.Equal(t, test.want.body, resBody, "Response body")
		})
	}
}

// Нужен для того, чтобы протестировать ошибку в http.ResponseWriter.Write()
type testResponseWriter struct {
	httptest.ResponseRecorder
	needError bool
}

func newTestResponseWriter(needError bool) *testResponseWriter {
	recoder := &testResponseWriter{needError: needError}
	recoder.Body = new(bytes.Buffer)

	return recoder
}

func (trw *testResponseWriter) Write(buf []byte) (int, error) {
	if trw.needError {
		return 0, fmt.Errorf("unable to write, for test only")
	}
	return trw.Body.Write(buf)
}
