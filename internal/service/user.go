package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repository.UserRepository
	hmacSecret []byte
}

func NewAuthService(repo repository.UserRepository, secret []byte) *AuthService {
	return &AuthService{
		repo:       repo,
		hmacSecret: secret,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input domain.SingUpInput) (int, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("service: hash password: %w", err)
	}

	user := &domain.User{
		Email:        input.Email,
		Name:         input.Name,
		Password:     string(hashed),
		RegisteredAt: time.Now(),
	}

	id, err := s.repo.CreateUser(ctx, *user)
	if err != nil {
		return 0, fmt.Errorf("service: create user: %w", err)
	}
	return id, nil
}

func (s *AuthService) SignIn(ctx context.Context, input domain.SingInInput) (string, error) {

	user, err := s.repo.GetByCredentials(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("service: get user by credentials: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", fmt.Errorf("service: incorrect password: %w", err)
	}

	// Используем RegisteredClaims — корректные имена полей и форматы
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Issuer:    "books-api", // опционально
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return "", fmt.Errorf("service: sign token: %w", err)
	}

	return signed, nil
}
