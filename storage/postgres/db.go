package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// сделать создание таблицы с помощью миграции sql файлами
func InitDB() (*sql.DB, error) {
	connStr := "postgres://postgres:76884@localhost:5433/postgres?sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
