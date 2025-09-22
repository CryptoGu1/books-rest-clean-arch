package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ConnectionInfo struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
	SSLMode  string
}

func NewPostgresConnectionInfo(info ConnectionInfo) (*sqlx.DB, error) {
	data := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		info.Host, info.Port, info.Username, info.DBName, info.Password, info.SSLMode)

	db, err := sqlx.Open("postgres", data)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %v", err)

	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return db, nil
}
