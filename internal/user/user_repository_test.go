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
		UserSettingId: 1,
	}
}
func mockUserRows(user User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "first_name", "last_name", "email"}).
		AddRow(user.Id, user.UserName, user.FirstName, user.LastName, user.Email)
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	userRepo := NewUserRepository(db)

	t.Run("Users Exists", func(t *testing.T) {
		rows := mockUserRows(mockUserObj())
		mock.ExpectQuery(sqlGetAllUsers).WillReturnRows(rows)

		users, err := userRepo.GetAllUsers()

		require.Nil(t, err)
		require.Equal(t, mockUserObj(), *users[0])
	})
}
