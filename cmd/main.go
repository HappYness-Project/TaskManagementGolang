package main

import (
	"fmt"
	"log"
	"os"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/cmd/configs"
	"github.com/happYness-Project/taskManagementGolang/cmd/db"
)

func main() {
	// how to get the environment from the application side?
	// fmt.Println("delve", isdelve.Enabled)
	var current_env = os.Getenv("APP_ENV")
	fmt.Println("Current Environment : " + current_env)
	env := configs.InitConfig(current_env)
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
