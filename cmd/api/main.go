package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"example.com/taskapp/internal/taskcontainer"
	"example.com/taskapp/internal/user"
)

const port = 8080

type application struct {
	DSN           string
	Domain        string
	containerRepo taskcontainer.ContainerRepository
	userRepo      user.UserRepository
	database      *sql.DB
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
	app.containerRepo = taskcontainer.NewContainerRepository(database)
	// app.userRepo = user.NewUserRepository(database)
	defer app.database.Close()

	app.Domain = "example.com"
	// err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	server := NewApiServer(fmt.Sprintf(":%d", port), database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
