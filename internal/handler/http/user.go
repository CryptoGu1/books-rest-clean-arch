package http

import (
	"net/http"
	"time"

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
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid input body"))
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
		"id": id,
	})

}

func (h *Handler) signIn(c echo.Context) error {
	var input domain.SingInInput

	if err := c.Bind(&input); err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-in",
			"problem": "bind error",
		}).Error(err)

		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid request body"))
	}

	if err := h.validate.Struct(input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	ctx := c.Request().Context()
	token, refresh, err := h.UserService.SignIn(ctx, input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "sign-in",
			"problem": "service error",
		}).Error(err)
		return respondErr(c, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh-token"
	cookie.Value = refresh
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30)

	c.SetCookie(cookie)

	return respondJSON(c, http.StatusOK, map[string]interface{}{
		"token": token,
	})

}

func (h *Handler) refresh(c echo.Context) error {
	// 1. Достаем cookie с refresh токеном
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token"))
	}

	refreshToken := cookie.Value
	if refreshToken == "" {
		return respondErr(c, echo.NewHTTPError(http.StatusUnauthorized, "empty refresh token"))
	}

	// 2. Обновляем токены через сервис
	ctx := c.Request().Context()
	accessToken, newRefreshToken, err := h.UserService.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusUnauthorized, err.Error()))
	}

	// 3. Перезаписываем cookie с новым refresh токеном
	newCookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // в продакшене поставить true
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	c.SetCookie(newCookie)

	// 4. Возвращаем новый access токен
	return respondJSON(c, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token"))
	}

	ctx := c.Request().Context()
	accessToken, refreshToken, err := h.UserService.RefreshTokens(ctx, cookie.Value)
	if err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusUnauthorized, err.Error()))
	}

	cookieSet := &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true в продакшене (если HTTPS)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	}
	c.SetCookie(cookieSet)

	return respondJSON(c, http.StatusOK, map[string]interface{}{
		"token": accessToken,
	})

}
