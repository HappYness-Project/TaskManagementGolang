package taskcontainer

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mockContainerObj() TaskContainer {
	return TaskContainer{
		ContainerId:   uuid.NewString(),
		ContainerName: "testuser",
		ContainerDesc: "testdesc",
	}
}
func mockContainerRows(container TaskContainer) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(container.ContainerId, container.ContainerName, container.ContainerDesc)
}

func TestContainerRepo_AllTaskContainers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	containerRepo := NewContainerRepository(db)

	t.Run("Container Exists", func(t *testing.T) {
		mockContainer := mockContainerObj()
		rows := mockContainerRows(mockContainer)
		mock.ExpectQuery(sqlGetAllContainers).WillReturnRows(rows)

		containers, err := containerRepo.AllTaskContainers()

		require.Nil(t, err)
		require.Equal(t, mockContainer, *containers[0])
	})

	t.Run("Container does not exist, then should return nothing with no result error msg", func(t *testing.T) {
		mock.ExpectQuery(sqlGetAllContainers).WillReturnError(sql.ErrNoRows)

		containers, err := containerRepo.AllTaskContainers()

		require.NotNil(t, err)
		require.Equal(t, "sql: no rows in result set", err.Error())
		require.Len(t, containers, 0)
	})
}

func TestContainerRepo_ContainerById(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	containerRepo := NewContainerRepository(db)
	t.Run("Container Exists", func(t *testing.T) {
		mockContainer := mockContainerObj()
		rows := mockContainerRows(mockContainer)
		mock.ExpectQuery("SELECT (.+) FROM public.taskcontainer WHERE id = \\$1").
			WithArgs(mockContainer.ContainerId).
			WillReturnRows(rows)

		container, err := containerRepo.GetById(mockContainer.ContainerId)

		require.Nil(t, err)
		require.Equal(t, &mockContainer, container)
	})

	t.Run("When Cannot find Container, Then return null", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM public.taskcontainer WHERE id = \\$1").
			WithArgs("invalid")

		container, err := containerRepo.GetById("invalid")

		require.NotNil(t, err)
		require.Nil(t, container)
	})

	t.Run("Internal server error.", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM public.taskcontainer WHERE id = \\$1").
			WithArgs("SomeError").WillReturnError(errors.New("random internal error occurred"))

		container, err := containerRepo.GetById("invalid")
		require.NotNil(t, err)
		require.Nil(t, container)
		// TODO Need to update GetById to handle different scenario.
	})
}
