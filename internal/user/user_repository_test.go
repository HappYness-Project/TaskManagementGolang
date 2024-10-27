package user

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func mockUserObj() User {
	return User{
		Id:            1,
		UserName:      "testuser",
		FirstName:     "kevin",
		LastName:      "park",
		Email:         "testuser@hproject.com",
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now().Add(time.Duration(20)),
		UserSettingId: 1,
	}
}
func mockUserRows(user User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "email", "is_active", "created_at", "updated_at", "usersetting_id"}).
		AddRow(user.Id, user.UserName, user.FirstName, user.LastName, user.Email, user.IsActive, user.CreatedAt, user.UpdatedAt, user.UserSettingId)
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()
	userRepo := NewUserRepository(db)

	t.Run("Users Exists", func(t *testing.T) {
		mockUser := mockUserObj()
		rows := mockUserRows(mockUser)
		mock.ExpectQuery(sqlGetAllUsers).WillReturnRows(rows)

		users, err := userRepo.GetAllUsers()

		require.Nil(t, err)
		require.Equal(t, mockUser, *users[0])
	})
}
