package tests

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/stretchr/testify/require"
)

var db *sql.DB

func TestMain(m *testing.M) {
	env := configs.InitConfig("")
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	logger := loggers.Setup(env)
	logger.Info().Msg(connStr)
	db, _ = dbs.ConnectToDb(connStr)

	m.Run()
}

func Test_TaskRepository_GetAllTasks_ReturnSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	taskRepo := repository.NewTaskRepository(db)

	tasks, err := taskRepo.GetAllTasks()

	require.Nil(t, err)
	require.NotNil(t, tasks)
}

func Test_TaskRepository_GetTaskById_ReturnSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	taskRepo := repository.NewTaskRepository(db)
	tasks, err := taskRepo.GetAllTasks()
	if err != nil {
		return
	}

	task, err := taskRepo.GetTaskById(tasks[0].TaskId)

	require.Nil(t, err)
	require.NotNil(t, task)
}
