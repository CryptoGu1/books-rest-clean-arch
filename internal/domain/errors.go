package domain

import "errors"

var (
	ErrBookNotFound         = errors.New("book not found")
	ErrRefreshTokenNotFound = errors.New("refresh token expired")
)
