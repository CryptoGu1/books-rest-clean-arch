package http

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/labstack/echo/v4"
)

func respondJSON(c echo.Context, code int, payload interface{}) error {
	return c.JSON(code, payload)
}

func respondErr(c echo.Context, err error) error {
	status := mapErrorToStatus(err)
	//На будущее есть ode, trace_id и т.д
	return c.JSON(status, map[string]interface{}{
		"error": err.Error(),
		"time":  time.Now().UTC(),
	})
}

func mapErrorToStatus(err error) int {
	var he *echo.HTTPError
	// domain.ErrBookNotFound -> 404
	if errors.Is(err, domain.ErrBookNotFound) {
		return http.StatusNotFound
	}
	// sql.ErrNoRows -> 404
	if errors.Is(err, sql.ErrNoRows) {
		return http.StatusNotFound
	}
	// по умолчанию — 500
	if errors.As(err, &he) {
		return he.Code
	}
	return http.StatusInternalServerError
}
