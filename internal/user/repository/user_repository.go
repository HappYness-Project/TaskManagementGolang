package repository

import (
	"database/sql"
	"fmt"

	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
)

type UserRepository interface {
	GetAllUsers() ([]*model.User, error)
	GetUserByUserId(userId string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUsersByGroupId(groupId int) ([]*model.User, error)
	CreateUser(user model.User) error
	UpdateUser(user model.User) error
}
type UserRepo struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (s *UserRepo) GetAllUsers() ([]*model.User, error) {
	rows, err := s.DB.Query(sqlGetAllUsers)
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0)
	for rows.Next() {
		p, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, p)
	}

	return users, nil
}

func (m *UserRepo) GetUserByUserId(user_id string) (*model.User, error) {
	rows, err := m.DB.Query(sqlGetUserByUserId, user_id)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (m *UserRepo) GetUserByEmail(email string) (*model.User, error) {
	rows, err := m.DB.Query(sqlGetUserByEmail, email)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}
func (m *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	rows, err := m.DB.Query(sqlGetUserByUsername, username)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (m *UserRepo) GetUsersByGroupId(groupId int) ([]*model.User, error) {
	rows, err := m.DB.Query(sqlGetUsersByGroupId, groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (m *UserRepo) CreateUser(user model.User) error {

	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction failure")
	}

	_, err = tx.Exec(sqlCreateUser, user.UserId, user.UserName, user.FirstName, user.LastName, user.Email, user.IsActive, user.CreatedAt, user.UpdatedAt, user.DefaultGroupId)
	if err != nil {
		return fmt.Errorf("unable to insert into user table : %w", err)
	}

	_, err = tx.Exec(sqlCreateUserSetting, user.Id, 0)
	if err != nil {
		return fmt.Errorf("unable to insert into usersetting table : %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit failure: %w", err)
	}

	return nil
}
func (m *UserRepo) UpdateUser(user model.User) error {
	_, err := m.DB.Exec(sqlUpdateUser, user.Id, user.FirstName, user.LastName, user.Email, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func scanRowsIntoUser(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.Id,
		&user.UserId,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DefaultGroupId,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
