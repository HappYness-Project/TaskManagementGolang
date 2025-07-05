package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/stretchr/testify/require"
)

func mockUserObj() model.User {
	return model.User{
		Id:             1,
		UserId:         "01960938-573a-723d-a4ac-cf3bbe420ece",
		UserName:       "testuser",
		FirstName:      "kevin",
		LastName:       "park",
		Email:          "testuser@hproject.com",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now().Add(time.Duration(20)),
		DefaultGroupId: 1,
	}
}
func mockUserRows(user model.User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "user_id", "username", "first_name", "last_name", "email", "is_active", "created_at", "updated_at", "default_group_id"}).
		AddRow(user.Id, user.UserId, user.UserName, user.FirstName, user.LastName, user.Email, user.IsActive, user.CreatedAt, user.UpdatedAt, user.DefaultGroupId)
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	env := configs.InitConfig("")
	logger := loggers.Setup(env)
	userRepo := NewUserRepository(db, logger)

	t.Run("Given single user, When Get All user, return single user", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)
		mock.ExpectQuery(sqlGetAllUsers).WillReturnRows(rows)

		users, err := userRepo.GetAllUsers()

		require.Nil(t, err)
		require.Equal(t, mockUser, *users[0])
	})
}
