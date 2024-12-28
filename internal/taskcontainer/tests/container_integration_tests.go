package taskcontainer

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/stretchr/testify/require"
)

func TestTaskContainerHandlerChecking(t *testing.T) {
	t.Parallel()
	env := configs.InitConfig("development.env")
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	database, _ := dbs.ConnectToDb(connStr)
	containerRepo := taskcontainer.NewContainerRepository(database)

	t.Run("Container Exists", func(t *testing.T) {
		container, err := containerRepo.AllTaskContainers()

		require.Nil(t, err)
		require.NotNil(t, container)
		t.Errorf("expected status code %d, got %d", http.StatusOK, 200)

	})
}
