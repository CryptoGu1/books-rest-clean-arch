package repository

import (
	"context"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Tokens struct {
	db *sqlx.DB
}

func NewToken(db *sqlx.DB) *Tokens {
	return &Tokens{db: db}
}

func (r *Tokens) Create(ctx context.Context, token domain.RefreshSession) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at)
		 VALUES ($1, $2, $3)`,
		token.UserID, token.Token, token.ExpiresAt)
	return err
}

func (r *Tokens) Get(ctx context.Context, token string) (domain.RefreshSession, error) {
	var t domain.RefreshSession
	query := `select user_id, token, expires_at from refresh_tokens where token = $1`
	err := r.db.GetContext(ctx, &t, query, token)
	return t, err
}
