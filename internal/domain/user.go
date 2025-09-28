package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

type SingUpInput struct {
	Name     string `json:"name" validate:"required,gte=2,lte=20""`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

func (i SingUpInput) Validate() error {
	return validate.Struct(i)
}

type SingInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

func (i SingInInput) Validate() error {
	return validate.Struct(i)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
