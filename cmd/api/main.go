package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
)

const port = 8080

type application struct {
	DSN      string
	Domain   string
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
	// err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	server := NewApiServer(fmt.Sprintf(":%d", port), database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
