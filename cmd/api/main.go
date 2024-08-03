package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"example.com/taskapp/internal/repository"
	"example.com/taskapp/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.DatabaseRepo
}

func main() {
	// Set Application Config
	var app application

	// Read from command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=8010 user=postgres password=postgres dbname=postgres sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	// connect to the db.
	conn, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDbRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.Domain = "example.com"
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
