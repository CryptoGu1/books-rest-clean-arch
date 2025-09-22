package main

import (
	"fmt"
	"log"

	_ "github.com/CryptoGu1/books-rest-clean-arch/docs"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/config"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/repository"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/service"
	"github.com/CryptoGu1/books-rest-clean-arch/internal/transport/http"
	"github.com/CryptoGu1/books-rest-clean-arch/pkg/postgres"
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

func main() {
	//init db
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v\n", cfg)

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
	log.Println("Connected to database")
	defer db.Close()

	//init DI
	bookRepo := repository.NewBookPostgresRepo(db)
	bookService := service.NewBookService(bookRepo)
	handler := http.NewHandler(bookService)

	router := handler.InitRouter()

	log.Println("Listening on port 8080")
	if err := router.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}

}
