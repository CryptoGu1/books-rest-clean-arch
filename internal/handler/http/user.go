package http

import (
	"net/http"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c echo.Context) error {
	var input domain.SingUpInput

	if err := c.Bind(&input); err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-up",
			"problem": "bind error",
		}).Error(err)

		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid request body"))
	}

	if err := h.validate.Struct(input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	ctx := c.Request().Context()
	id, err := h.UserService.SignUp(ctx, input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-up",
			"problem": "service error",
		}).Error(err)
		return respondErr(c, err)

	}

	return respondJSON(c, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "user created successfully",
	})

}

func (h *Handler) signIn(c echo.Context) error {
	var input domain.SingInInput

	if err := c.Bind(&input); err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-up",
			"problem": "bind error",
		}).Error(err)

		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid request body"))
	}

	if err := h.validate.Struct(input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	ctx := c.Request().Context()
	token, err := h.UserService.SignIn(ctx, input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-in",
			"problem": "service error",
		}).Error(err)
		return respondErr(c, err)
	}
	return respondJSON(c, http.StatusOK, map[string]interface{}{
		"token": token,
	})

}
