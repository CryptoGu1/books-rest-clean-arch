package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		logrus.WithFields(logrus.Fields{
			"method": req.Method,
			"uri":    req.RequestURI,
			"remote": req.RemoteAddr,
		}).Info("incoming request")

		return next(c)
	}
}

func (h *Handler) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")

		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")

		}
		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return h.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid claims")

		}

		userId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid subject")

		}

		c.Set("userID", userId)
		return next(c)
	}
}
