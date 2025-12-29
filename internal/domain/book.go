package domain

import (
	"encoding/json"
	"time"
)

// --- Кастомный тип ---
type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	// убираем кавычки
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// если пустая строка → пропускаем
	if s == "" {
		return nil
	}

	// парсим как "YYYY-MM-DD"
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	// сохраняем в UTC
	*d = DateOnly(t.UTC())
	return nil
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}

// --- Модели ---
type Book struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Author      string    `db:"author" json:"author"`
	PublishDate time.Time `db:"publish_date" json:"publish_date"`
	Rating      int       `db:"rating" json:"rating"`
}

type CreateBookInput struct {
	Title       string    `json:"title" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	PublishDate *DateOnly `json:"publish_date,omitempty"`
	Rating      int       `json:"rating" validate:"min=1,max=5"`
}

type UpdateBookInput struct {
	Title       string    `json:"title" validate:"required"`
	Author      string    `json:"author" validate:"required"`
	PublishDate *DateOnly `json:"publish_date,omitempty"`
	Rating      int       `json:"rating" validate:"min=1,max=5"`
}

// --- Конвертации ---
func (input *CreateBookInput) ToBook() *Book {
	book := &Book{
		Title:  input.Title,
		Author: input.Author,
		Rating: input.Rating,
	}

	if input.PublishDate != nil {
		book.PublishDate = input.PublishDate.ToTime()
	} else {
		book.PublishDate = time.Now().UTC()
	}
	return book
}

func (input *UpdateBookInput) ToBook() *Book {
	book := &Book{
		Title:  input.Title,
		Author: input.Author,
		Rating: input.Rating,
	}

	if input.PublishDate != nil {
		book.PublishDate = input.PublishDate.ToTime()
	} else {
		book.PublishDate = time.Now().UTC()
	}
	return book
}
