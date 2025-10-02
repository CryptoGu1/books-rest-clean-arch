package http

import (
	"net/http"
	"strconv"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// Create godoc
//
//	@Summary		Create new book
//	@Description	Добавляет новую книгу в базу данных books
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			input	body		domain.CreateBookInput	true	"title, author, publish_date(default null), rating"
//	@Success		201		{object}	domain.Book
//	@Failure		400		{object}	map[string]string
//	@Router			/books [post]
func (h *Handler) Create(c echo.Context) error {
	var input domain.CreateBookInput
	if err := c.Bind(&input); err != nil {
		log.WithFields(log.Fields{
			"handler": "CreateBook",
			"problem": "request body bind error",
		}).Error(err)

		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid request body"))
	}

	if err := h.validate.Struct(input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	ctx := c.Request().Context()
	id, err := h.bookService.Create(ctx, &input)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "CreateBook",
			"problem": "service error",
		}).Error(err)

		return respondErr(c, err)
	}

	return respondJSON(c, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "Book created successfully",
	})
}

// GetById godoc
//
//	@Summary		Get Book by id
//	@Description	получает книгу по id
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int true "Book ID"
//	@Success		200		{object}	domain.Book
//	@Failure		400		{object}	map[string]string
//	@Router			/books/{id} [get]
func (h *Handler) GetById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid book id"))

	}

	ctx := c.Request().Context()
	book, err := h.bookService.GetById(ctx, id)
	if err != nil {
		return respondErr(c, err)
	}

	return respondJSON(c, http.StatusOK, book)

}

// GetAllBooks godoc
// @Summary      Get all books
// @Description  Возвращает список всех книг
// @Tags         books
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Book
// @Failure      500  {object}  map[string]string
// @Router       /books [get]
func (h *Handler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	books, err := h.bookService.GetAll(ctx)
	if err != nil {
		return respondErr(c, err)

	}
	return respondJSON(c, http.StatusOK, books)
}

// UpdateBook godoc
// @Summary      Update book
// @Description  Обновляет данные книги по ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id     path      int                  true  "Book ID"
// @Param        input  body      domain.UpdateBookInput  true  "Updated book data"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /books/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid book id"))
	}

	var input domain.UpdateBookInput
	if err := c.Bind(&input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid request body"))
	}
	if err := h.validate.Struct(&input); err != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, err.Error()))

	}

	book := input.ToBook()
	ctx := c.Request().Context()
	if err := h.bookService.Update(ctx, id, book); err != nil {
		return respondErr(c, err)
	}
	return respondJSON(c, http.StatusOK, map[string]interface{}{
		"message": "Book updated successfully",
		"id":      id,
	})

}

// DeleteBook godoc
// @Summary      Delete book
// @Description  Удаляет книгу по ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /books/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, errConv := strconv.Atoi(idParam)
	if errConv != nil {
		return respondErr(c, echo.NewHTTPError(http.StatusBadRequest, "invalid book id"))

	}

	ctx := c.Request().Context()
	if err := h.bookService.Delete(ctx, id); err != nil {
		return respondErr(c, err)
	}

	return c.NoContent(http.StatusNoContent)

}
