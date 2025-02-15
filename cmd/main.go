package main

import (
	"fmt"
	"os"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func main() {
	var current_env = os.Getenv("APP_ENV")
	fmt.Println("Current Environment : " + current_env)
	env := configs.InitConfig(current_env)
	logger := loggers.Setup(env)
	configs.AccessToken = env.AccessTokenSecret

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s timezone=UTC connect_timeout=5 ",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	if current_env == "development" {
		connStr += "sslmode=require"
	} else {
		connStr += "sslmode=disable"
	}

	logger.Info().Msg(connStr)
	database, err := dbs.ConnectToDb(connStr)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to connect to the database.")
		return
	}

	server := api.NewApiServer(fmt.Sprintf("%s:%d", env.Host, env.Port), database, logger)
	r := server.Setup()
	if err := server.Run(r); err != nil {
		logger.Error().Err(err).Msg("Unable to set up the server.")
		return
	}
}
