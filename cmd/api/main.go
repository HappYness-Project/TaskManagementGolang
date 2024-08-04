package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"example.com/taskapp/internal/repository"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.ContainerRepository
	// UserDBRepo repository.UserRepository
	database *sql.DB
}

func main() {
	// Set Application Config
	var app application

	// Read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=8010 user=postgres password=postgres dbname=postgres sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	// connect to the db.
	database, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &repository.ContainerRepo{DB: database}
	// app.UserDBRepo = &dbrepo.
	defer app.database.Close()

	app.Domain = "example.com"
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
