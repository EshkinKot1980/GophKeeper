package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

const (
	Scheme       = "https://"
	APIprefix    = "/api"
	RegisterPath = "/register"
	LoginPath    = "/login"
	SecretPath   = "/secret"
	ContentType  = "application/json"
)

var (
	ErrRegistrationFailed   = errors.New("failed to register user")
	ErrLoginFailed          = errors.New("login failed")
	ErrSecretSendFailed     = errors.New("failed to send secret")
	ErrSecretRetrieveFailed = errors.New("failed to retrieve secret")
)

type Client struct {
	baseURL string
	client  *resty.Client
}

func NewClient(baseURL string, allowSefSignedCert bool) *Client {
	c := Client{
		baseURL: baseURL,
		client: resty.New().
			SetTimeout(time.Minute).
			SetBaseURL(baseURL).
			SetHeader("Content-Type", ContentType),
	}

	if allowSefSignedCert {
		c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &c
}

// Register регистрирует пользователя в системе.
func (c *Client) Register(cr dto.Credentials) (dto.AuthResponse, error) {
	var authResp dto.AuthResponse
	req := c.client.R().
		SetResult(&authResp).
		SetBody(cr)

	resp, err := req.Post(RegisterPath)
	if err != nil {
		return authResp, fmt.Errorf("%w: %w", ErrRegistrationFailed, err)
	} else if !resp.IsSuccess() {
		return authResp, fmt.Errorf("%w: %s", ErrRegistrationFailed, resp.String())
	}

	return authResp, nil
}

// Login осуществляет вход пользователя в систему.
func (c *Client) Login(cr dto.Credentials) (dto.AuthResponse, error) {
	var authResp dto.AuthResponse

	req := c.client.R().
		SetResult(&authResp).
		SetBody(cr)

	resp, err := req.Post(LoginPath)
	if err != nil {
		return authResp, fmt.Errorf("%w: %w", ErrLoginFailed, err)
	} else if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusInternalServerError {
			return authResp, fmt.Errorf("%w: internal server error", ErrLoginFailed)
		}

		return authResp, ErrLoginFailed
	}

	return authResp, nil
}

// Upload coхраняет секрет на сервере.
func (c *Client) Upload(data dto.SecretRequest, token string) error {
	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetBody(data)

	resp, err := req.Post(SecretPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSecretSendFailed, err)
	} else if !resp.IsSuccess() {
		text := http.StatusText(resp.StatusCode())
		if body := resp.String(); body != "" {
			text += ": " + body
		}

		return fmt.Errorf("%w: %s", ErrSecretSendFailed, text)
	}

	return nil
}

// Retrieve получает секрет пользователя с ервера
func (c *Client) Retrieve(id uint64, token string) (dto.SecretResponse, error) {
	var secret dto.SecretResponse

	req := c.client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&secret)

	path := fmt.Sprintf("%s/%d", SecretPath, id)
	resp, err := req.Get(path)

	if err != nil {
		return secret, fmt.Errorf("%w: %w", ErrSecretRetrieveFailed, err)
	} else if !resp.IsSuccess() {
		switch resp.StatusCode() {
		case http.StatusUnauthorized:
			return secret, fmt.Errorf("%w: authorization failed", ErrSecretRetrieveFailed)
		case http.StatusBadRequest:
			return secret, fmt.Errorf("%w: %s", ErrSecretRetrieveFailed, resp)
		case http.StatusNotFound:
			return secret, fmt.Errorf("%w: not found", ErrSecretRetrieveFailed)
		default:
			return secret, fmt.Errorf("%w: internal server error", ErrSecretRetrieveFailed)
		}
	}

	return secret, nil
}
