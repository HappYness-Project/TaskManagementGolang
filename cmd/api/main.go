package main

import (
	"flag"
	"fmt"
	"log"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
}

func main() {
	var app application

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=8010 user=postgres password=postgres dbname=postgres sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	database, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}
	server := NewApiServer(fmt.Sprintf(":%d", port), database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
