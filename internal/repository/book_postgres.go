package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) (int, error)
	GetBook(ctx context.Context, id int) (*domain.Book, error)
	GetAllBooks(ctx context.Context) ([]*domain.Book, error)
	Update(ctx context.Context, id int, book *domain.Book) error
	Delete(ctx context.Context, id int) error
}

type BookPostgresRepo struct {
	db *sqlx.DB
}

func NewBookPostgresRepo(db *sqlx.DB) BookRepository {
	return &BookPostgresRepo{db: db}

}

func (r *BookPostgresRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repo: delete book: %w", err)
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: delete rows affected: %w", err)
	}
	if aff == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}

func (r *BookPostgresRepo) Create(ctx context.Context, book *domain.Book) (int, error) {
	var id int
	query := `
	INSERT INTO books (title, author, publish_date, rating)
	values ($1, $2, $3, $4) RETURNING id`
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	err := r.db.QueryRowxContext(ctx, query, book.Title, book.Author, book.PublishDate, book.Rating).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: create book: %w", err)
	}

	book.ID = id
	return id, nil
}

func (r *BookPostgresRepo) GetBook(ctx context.Context, id int) (*domain.Book, error) {
	var book domain.Book
	query := `
	SELECT * FROM books WHERE id = $1`

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBookNotFound
		}
		return nil, fmt.Errorf("repo: get book: %w", err)
	}
	return &book, nil
}

func (r *BookPostgresRepo) GetAllBooks(ctx context.Context) ([]*domain.Book, error) {
	var books []*domain.Book
	query := `
	SELECT * FROM books`

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	if err := r.db.SelectContext(ctx, &books, query); err != nil {
		return nil, fmt.Errorf("repo: get all books: %w", err)
	}

	return books, nil
}

func (r *BookPostgresRepo) Update(ctx context.Context, id int, book *domain.Book) error {
	query := `UPDATE books SET title=$1, author=$2, publish_date=$3, rating=$4 WHERE id=$5`
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	res, err := r.db.ExecContext(ctx, query, book.Title, book.Author, book.PublishDate, book.Rating, id)
	if err != nil {
		return fmt.Errorf("repo: update book: %w", err)

	}

	aff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: update rows affected: %w", err)
	}
	if aff == 0 {
		return domain.ErrBookNotFound
	}
	return nil
}
