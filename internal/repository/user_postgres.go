package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, input domain.User) (int, error)
	GetByCredentials(ctx context.Context, email string) (domain.User, error)
}
type UserPostgresRepo struct {
	db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) CreateUser(ctx context.Context, input domain.User) (int, error) {
	var id int
	query := `insert into users (name, email, password_hash, registered_at) values ($1, $2, $3, $4) returning id`

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	err := r.db.QueryRowxContext(ctx, query, input.Name, input.Email, input.Password, input.RegisteredAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo:error creating user: %w", err)
	}
	return id, nil
}

func (r *UserPostgresRepo) GetByCredentials(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := `select * from users where email = $1 `

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	err := r.db.QueryRowxContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RegisteredAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("repo:error getting user by credentials: %w", err)
	}

	return user, err
}
