package http

import (
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
