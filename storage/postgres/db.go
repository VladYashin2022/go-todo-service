package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// сделать создание таблицы с помощью миграции sql файлами
func InitDB() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://postgres:76884@localhost:5433/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {

		if err := db.Ping(); err == nil {
			return db, nil
		}
		log.Printf("waiting for database... attempt %d", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("Database is not available")
}