package service

import (
	"context"
	"fmt"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/repository"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService {
	return &BookService{repo}
}

func (s *BookService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: delete book : %w", err)

	}
	return nil
}

func (s *BookService) Create(ctx context.Context, input *domain.CreateBookInput) (int, error) {
	book := input.ToBook()
	id, err := s.repo.Create(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("service: create book: %w", err)
	}
	return id, nil
}

func (s *BookService) GetById(ctx context.Context, id int) (*domain.Book, error) {
	book, err := s.repo.GetBook(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: get book: %w", err)

	}

	return book, nil
}

func (s *BookService) GetAll(ctx context.Context) ([]*domain.Book, error) {
	books, err := s.repo.GetAllBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: get all books: %w", err)
	}
	return books, nil
}

func (s *BookService) Update(ctx context.Context, id int, book *domain.Book) error {
	err := s.repo.Update(ctx, id, book)
	if err != nil {
		return fmt.Errorf("service: update book: %w", err)
	}

	return nil

}

//func (s *BookService) validateCreateInput(input *domain.CreateBookInput) error {
//	if input.Title == "" {
//		return errors.New("title is required")
//	}
//
//	if input.Author == "" {
//		return errors.New("author is required")
//	}
//
//	if input.Rating < 0 || input.Rating > 5 {
//		return errors.New("rating must be between 0 and 5")
//	}
//
//	return nil
//}
