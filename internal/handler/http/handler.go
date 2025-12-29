package http

import (
	"context"

	"github.com/CryptoGu1/books-rest-clean-arch/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type BookService interface {
	Create(ctx context.Context, input *domain.CreateBookInput) (int, error)
	GetById(ctx context.Context, id int) (*domain.Book, error)
	GetAll(ctx context.Context) ([]*domain.Book, error)
	Update(ctx context.Context, id int, book *domain.Book) error
	Delete(ctx context.Context, id int) error
}

type AuthService interface {
	SignUp(ctx context.Context, input domain.SingUpInput) (int, error)
	SignIn(ctx context.Context, input domain.SingInInput) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	bookService BookService
	UserService AuthService
	validate    *validator.Validate
	jwtSecret   []byte
}

func NewHandler(bookService BookService, userService AuthService, jwtSecret []byte) *Handler {
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
		auth.POST("/refresh", h.refresh)

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
