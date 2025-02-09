package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println("Current Dir: " + exPath)

	var current_env = os.Getenv("APP_ENV")
	fmt.Println("Current Environment : " + current_env)
	env := configs.InitConfig(current_env)
	logger := loggers.Setup(env)
	configs.AccessToken = env.AccessTokenSecret

	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	logger.Info().Msg(connStr)
	database, err := dbs.ConnectToDb(connStr)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to connect to the database.")
		return
	}

	server := api.NewApiServer(fmt.Sprintf(":%d", env.Port), database, logger)
	r := server.Setup()
	if err := server.Run(r); err != nil {
		logger.Error().Err(err).Msg("Unable to set up the server.")
		return
	}
}
