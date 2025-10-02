package http

import (
	"github.com/CryptoGu1/books-rest-clean-arch/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	bookService *service.BookService
	UserService *service.AuthService
	validate    *validator.Validate
	jwtSecret   []byte
}

func NewHandler(bookService *service.BookService, userService *service.AuthService, jwtSecret []byte) *Handler {
	return &Handler{
		bookService: bookService,
		UserService: userService,
		validate:    validator.New(),
		jwtSecret:   jwtSecret,
	}
}

func (h *Handler) InitRouter() *echo.Echo {
	e := echo.New()

	//Middlewares
	e.Use(LoggingMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//Swager
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)

	}

	booksGroup := e.Group("/books")
	booksGroup.Use(h.JWTMiddleware)
	{
		booksGroup.POST("", h.Create)
		booksGroup.GET("/:id", h.GetById)
		booksGroup.GET("", h.GetAll)
		booksGroup.PUT("/:id", h.Update)
		booksGroup.DELETE("/:id", h.Delete)
	}

	return e
}
