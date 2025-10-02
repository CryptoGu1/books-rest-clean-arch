package main

import (
	"fmt"
	"os"

	_ "github.com/CryptoGu1/books-rest-clean-arch/docs"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/config"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/handler/http"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/repository"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/service"
	"github.com/CryptoGu1/books-rest-clean-arch/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

//	@title			Swagger Books api
//	@version		1.0
//	@description	This is a simple RESTful api crud using Echo Framework
//	@host			localhost:8080
//	@BasePath		/

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	//init db
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresConnectionInfo(postgres.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		SSLMode:  cfg.DB.SSLMode,
		DBName:   cfg.DB.Name,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	//init DI
	bookRepo := repository.NewBookPostgresRepo(db)
	userRepo := repository.NewUserPostgresRepo(db)

	bookService := service.NewBookService(bookRepo)
	userService := service.NewAuthService(userRepo, jwtSecret)

	handler := http.NewHandler(bookService, userService, jwtSecret)

	router := handler.InitRouter()

	log.Info("SERVER STARTED")
	if err := router.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}

}
