package user

import (
	"database/sql"
)

type UserRepository interface {
	GetAllUsers() ([]*User, error)
	GetUserById(id int) (*User, error)
	// GetUsersByGroupId(groupid int) ([]*models.User, error)
	// Create(user *models.User) error
	// Update(user *models.User) error
}
type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (s *UserRepo) GetAllUsers() ([]*User, error) {
	rows, err := s.DB.Query("SELECT * FROM public.user")
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for rows.Next() {
		p, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, p)
	}

	return users, nil
}
func scanRowsIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)

	err := rows.Scan(
		&user.Id,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (m *UserRepo) GetUserById(id int) (*User, error) {
	var user *User
	err := m.DB.QueryRow("select * from public.user where id = ?", id).Scan(id)
	return user, err
}
