package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lechitz/AionApi/src/config"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.StringConnection)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
