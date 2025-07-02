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
	configs.AccessToken = "71871847e4548334f720bf055f30829e28f58a52bb4aae7319d5d775622682cf6ba54671a2c270110be13ffb3fea16b3563e2109a4d24612ac5c5469d9cbc9e5"

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s timezone=UTC connect_timeout=5 ",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	if current_env == "local" || current_env == "" {
		connStr += "sslmode=disable"
	} else {
		connStr += "sslmode=require"
	}

	logger.Info().Msg(connStr)
	database, err := dbs.ConnectToDb(connStr)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to connect to the database.")
		return
	}

	server := api.NewApiServer(fmt.Sprintf("%s:%s", env.Host, env.Port), database, logger)
	r := server.Setup()
	if err := server.Run(r); err != nil {
		logger.Error().Err(err).Msg("Unable to set up the server.")
		return
	}
}
