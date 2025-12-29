package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	audit "github.com/CryptoGu1/books-grpc-log/pkg/domain"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type SessionRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}
type AuthService struct {
	repo        repository.UserRepository
	sessionRepo SessionRepository
	auditClient AuditClient
	hmacSecret  []byte
}

func NewAuthService(repo repository.UserRepository, sessionRepo SessionRepository, auditClient AuditClient, secret []byte) *AuthService {
	return &AuthService{
		repo:        repo,
		sessionRepo: sessionRepo,
		auditClient: auditClient,
		hmacSecret:  secret,
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

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_REGISTER,
		Entity:    audit.ENTITY_USER,
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "SignUp",
		}).Error("failed to send log request", err)
	}
	return id, nil
}

func (s *AuthService) SignIn(ctx context.Context, input domain.SingInInput) (string, string, error) {

	user, err := s.repo.GetByCredentials(ctx, input.Email)
	if err != nil {
		return "", "", fmt.Errorf("service: get user by credentials: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", "", fmt.Errorf("service: incorrect password: %w", err)
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_LOGIN,
		Entity:    audit.ENTITY_USER,
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "SignIn",
		}).Error("failed to send log request", err)
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *AuthService) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	// Используем RegisteredClaims — корректные имена полей и форматы
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", fmt.Errorf("service: sign token: %w", err)
	}

	refresh, err := s.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("service: refresh token: %w", err)
	}

	if err := s.sessionRepo.Create(ctx, domain.RefreshSession{
		UserID:    userId,
		Token:     refresh,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", fmt.Errorf("service: create refresh token: %w", err)
	}

	return accessToken, refresh, nil
}

func (s *AuthService) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("service: refresh token: %w", err)
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenNotFound
	}

	return s.generateTokens(ctx, session.UserID)

}
