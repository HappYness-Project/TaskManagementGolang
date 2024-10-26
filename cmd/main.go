package main

import (
	"fmt"
	"log"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/cmd/configs"
	"github.com/happYness-Project/taskManagementGolang/cmd/db"
)

func main() {
	env := configs.InitConfig("development.env")
	configs.AccessToken = env.AccessTokenSecret
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)

	log.Print(connStr)
	database, err := db.ConnectToDb(connStr)
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewApiServer(fmt.Sprintf(":%d", env.Port), database)
	r := server.Setup()
	if err := server.Run(r); err != nil {
		log.Fatal(err)
	}
}
