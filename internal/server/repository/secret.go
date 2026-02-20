package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository/pg"
)

type Secret struct {
	pool *pgxpool.Pool
}

func NewSecret(db *pg.DB) *Secret {
	return &Secret{pool: db.Pool()}
}

// Create создает пользовательский секрет в БД.
func (s *Secret) Create(ctx context.Context, secret entity.Secret) error {
	query := `
	INSERT INTO secrets
			(user_id, data_type, name, meta_data, encrypted_data, encrypted_key) 
		VALUES
			($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at`

	_, err := s.pool.Exec(
		ctx,
		query,
		secret.UserID,
		secret.DataType,
		secret.Name,
		secret.MetaData,
		secret.EncryptedData,
		secret.EncryptedKey,
	)

	if err != nil {
		return fmt.Errorf("failed to insert to secrets: %w", errors.Trasform(err))
	}

	return nil
}

// GetForUser возвращает пользовательский секрет по secretID и userID.
func (s *Secret) GetForUser(ctx context.Context, secretID uint64, userID string) (entity.Secret, error) {
	var secret entity.Secret

	query := `SELECT * FROM secrets WHERE id = $1 AND user_id = $2`
	rows, err := s.pool.Query(ctx, query, secretID, userID)
	if err != nil {
		return secret, fmt.Errorf("failed to select from secrets: %w", err)
	}

	secret, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Secret])
	if err != nil {
		return secret, errors.Trasform(err)
	}

	return secret, nil
}
