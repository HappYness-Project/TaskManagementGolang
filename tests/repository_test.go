package tests

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/stretchr/testify/require"
)

var db *sql.DB

const dbTimeout = time.Second * 5

func TestMain(m *testing.M) {
	env := configs.InitConfig("")
	configs.AccessToken = env.AccessTokenSecret
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	logger := loggers.Setup(env)
	logger.Info().Msg(connStr)
	db, _ = dbs.ConnectToDb(connStr)

	m.Run()
}

func Test_TaskRepository_GetAllTasks_ReturnSuccess(t *testing.T) {
	taskRepo := repository.NewTaskRepository(db)

	tasks, err := taskRepo.GetAllTasks()

	require.Nil(t, err)
	require.NotNil(t, tasks)
}

// TODO : Require to test timeout scenario