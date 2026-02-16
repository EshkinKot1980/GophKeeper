package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository/pg"
)

type User struct {
	pool *pgxpool.Pool
}

func NewUser(db *pg.DB) *User {
	return &User{pool: db.Pool()}
}

func (u *User) GetByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, login, hash, auth_salt, encr_salt, created_at FROM users WHERE id = $1`
	row := u.pool.QueryRow(ctx, query, id)

	err := row.Scan(&user.ID, &user.Login, &user.Hash, &user.AuthSalt, &user.EncrSalt, &user.Created)
	if err != nil {
		return entity.User{}, errors.Trasform(err)
	}

	return user, nil
}

func (u *User) FindByLogin(ctx context.Context, login string) (entity.User, error) {
	var user entity.User
	query := `SELECT id, login, hash, auth_salt, encr_salt, created_at FROM users WHERE login = $1`
	row := u.pool.QueryRow(ctx, query, login)

	err := row.Scan(&user.ID, &user.Login, &user.Hash, &user.AuthSalt, &user.EncrSalt, &user.Created)
	if err != nil {
		return entity.User{}, errors.Trasform(err)
	}

	return user, nil
}

func (u *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := `INSERT INTO users (login, hash, auth_salt, encr_salt) VALUES($1, $2, $3, $4) RETURNING id, created_at`
	row := u.pool.QueryRow(ctx, query, user.Login, user.Hash, user.AuthSalt, user.EncrSalt)

	err := row.Scan(&user.ID, &user.Created)
	if err != nil {
		return user, fmt.Errorf("failed to insert to users: %w", errors.Trasform(err))
	}

	return user, nil
}
