package user

import (
	"database/sql"
)

type UserRepository interface {
	GetAllUsers() ([]*User, error)
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUsersByGroupId(groupId int) ([]*User, error)
	Create(user User) error
}
type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (s *UserRepo) GetAllUsers() ([]*User, error) {
	rows, err := s.DB.Query(sqlGetAllUsers)
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

func (m *UserRepo) GetUserById(id int) (*User, error) {
	var user *User
	err := m.DB.QueryRow(sqlGetUserById, id).Scan(id)
	return user, err
}

func (m *UserRepo) GetUserByEmail(email string) (*User, error) {
	var user *User
	err := m.DB.QueryRow(sqlGetUserByEmail, email).Scan(email)
	return user, err
}
func (m *UserRepo) GetUsersByGroupId(groupId int) ([]*User, error) {
	rows, err := m.DB.Query(sqlGetUsersByGroupId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (m *UserRepo) Create(user User) error {
	return nil
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
		&user.IsActive,
		&user.UserSettingId,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
