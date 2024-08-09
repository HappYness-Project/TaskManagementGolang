package main

import (
	"fmt"
	"log"

	"example.com/taskapp/cmd/configs"
)

func main() {
	env := configs.NewEnv()
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)

	log.Print(connStr)
	database, err := connectToDb(connStr)
	if err != nil {
		log.Fatal(err)
	}
	server := NewApiServer(fmt.Sprintf(":%d", env.Port), database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
