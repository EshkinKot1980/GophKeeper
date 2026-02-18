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
	ContentType  = "application/json"
)

var (
	ErrRegistrationFailed = errors.New("failed to register user")
	ErrLoginFailed        = errors.New("login failed")
)

type Client struct {
	baseURL string
	client  *resty.Client
}

func NewClient(serverAddr string, allowSefSignedCert bool) *Client {
	url := Scheme + serverAddr + APIprefix

	c := Client{
		baseURL: url,
		client: resty.New().
			SetTimeout(time.Minute).
			SetBaseURL(url).
			SetHeader("Content-Type", ContentType),
	}

	if allowSefSignedCert {
		c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &c
}

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
