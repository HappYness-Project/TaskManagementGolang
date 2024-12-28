package dbs

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDb(connStr string) (*sql.DB, error) {
	conn, err := openDb(connStr)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres")
	return conn, nil
}
