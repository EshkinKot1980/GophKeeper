// Пакет service содержит сервисный слой серверной части приложения
package service

import (
	"context"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	repErrors "github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByLogin(ctx context.Context, login string) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type Logger interface {
	Error(message string, err error)
}

// Auth сервис для регистрации аутентификации и авторизации
type Auth struct {
	repository UserRepository
	logger     Logger
	pub        *rsa.PublicKey
	priv       *rsa.PrivateKey
}

func NewAuth(r UserRepository, l Logger, jwtPub *rsa.PublicKey, jwtPriv *rsa.PrivateKey) *Auth {
	return &Auth{repository: r, logger: l, pub: jwtPub, priv: jwtPriv}
}

// Register регистрация пользователя по логину с паролем
func (a *Auth) Register(ctx context.Context, c dto.Credentials) (resp dto.AuthResponse, err error) {
	cr := trimCredentials(c)
	if err := cr.Validate(); err != nil {
		return resp, fmt.Errorf("%w: %w", srvErrors.ErrAuthInvalidCredentials, err)
	}

	authSalt, err := crypto.GenerateRandomBytes(crypto.SaltLen)
	if err != nil {
		a.logger.Error("failed to generate auth salt", err)
		return resp, srvErrors.ErrUnexpected
	}

	hash, err := crypto.DeriveKey([]byte(c.Password), authSalt)
	if err != nil {
		a.logger.Error("failed to hash password", err)
		return resp, srvErrors.ErrUnexpected
	}

	encrSalt, err := crypto.GenerateRandomBytes(crypto.SaltLen)
	if err != nil {
		a.logger.Error("failed to generate encriptin salt", err)
		return resp, srvErrors.ErrUnexpected
	}
	base64EcrSalt := base64.RawStdEncoding.EncodeToString(encrSalt)

	user := entity.User{
		Login:    c.Login,
		Hash:     base64.RawStdEncoding.EncodeToString(hash),
		AuthSalt: base64.RawStdEncoding.EncodeToString(authSalt),
		EncrSalt: base64EcrSalt,
	}
	user, err = a.repository.Create(ctx, user)

	if err != nil {
		switch {
		case errors.Is(err, repErrors.ErrDuplicateKey):
			return resp, srvErrors.ErrAuthUserAlreadyExists
		default:
			a.logger.Error("failed to create user", err)
			return resp, srvErrors.ErrUnexpected
		}
	}

	token, err := a.generateToken(user)
	if err != nil {
		return resp, err
	}

	resp.Token = token
	resp.EncrSalt = base64EcrSalt
	return resp, nil
}

// Login вход пользователя в систему по логину с паролем
func (a *Auth) Login(ctx context.Context, c dto.Credentials) (resp dto.AuthResponse, err error) {
	cr := trimCredentials(c)

	user, err := a.repository.FindByLogin(ctx, cr.Login)
	if err != nil {
		if errors.Is(err, repErrors.ErrNotFound) {
			return resp, srvErrors.ErrAuthInvalidCredentials
		} else {
			a.logger.Error("failed to find user", err)
			return resp, srvErrors.ErrUnexpected
		}
	}

	authSalt, err := base64.RawStdEncoding.DecodeString(user.AuthSalt)
	if err != nil {
		a.logger.Error("failed decode user auth salt", err)
		return resp, srvErrors.ErrUnexpected
	}

	userHash, err := base64.RawStdEncoding.DecodeString(user.Hash)
	if err != nil {
		a.logger.Error("failed decode user hash", err)
		return resp, srvErrors.ErrUnexpected
	}

	calculatedHash, err := crypto.DeriveKey([]byte(c.Password), authSalt)
	if err != nil {
		a.logger.Error("failed to hash password", err)
		return resp, srvErrors.ErrUnexpected
	}

	if subtle.ConstantTimeCompare(userHash, calculatedHash) != 1 {
		return resp, srvErrors.ErrAuthInvalidCredentials
	}

	token, err := a.generateToken(user)
	if err != nil {
		return resp, err
	}

	resp.Token = token
	resp.EncrSalt = user.EncrSalt
	return resp, nil
}

// User получение пользователя по JWT
func (a *Auth) User(ctx context.Context, token string) (entity.User, error) {
	var user entity.User

	jt, err := jwt.Parse(
		token,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				err := srvErrors.ErrAuthInvalidToken
				a.logger.Error("unexpected signing method: "+t.Method.Alg(), err)
				return nil, err
			}
			return a.pub, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return user, srvErrors.ErrAuthTokenExpired
		}
		return user, srvErrors.ErrAuthInvalidToken
	}

	userID, err := parseID(jt)
	if err != nil {
		return user, srvErrors.ErrAuthInvalidToken
	}

	user, err = a.repository.GetByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repErrors.ErrNotFound) {
			a.logger.Error("failed to find user by id", err)
		}
		return user, srvErrors.ErrAuthInvalidToken
	}

	return user, nil
}

func parseID(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", srvErrors.ErrAuthInvalidToken
	}

	id, ok := claims["jti"].(string)
	if !ok {
		return "", srvErrors.ErrAuthInvalidToken
	}

	return id, nil
}

func (a *Auth) generateToken(u entity.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			ID:        u.ID,
		},
	)

	tokenStr, err := token.SignedString(a.priv)
	if err != nil {
		a.logger.Error("failed to generate token", err)
		return "", srvErrors.ErrUnexpected
	}

	return tokenStr, nil
}

func trimCredentials(c dto.Credentials) dto.Credentials {
	return dto.Credentials{
		Login:    strings.TrimSpace(c.Login),
		Password: strings.TrimSpace(c.Password),
	}
}
