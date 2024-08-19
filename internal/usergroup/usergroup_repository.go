package usergroup

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 5

type UserGroupRepository interface {
	GetAllUsergroups() ([]*UserGroup, error)
	GetById(id int) (*UserGroup, error)
	GetUserGroupsByUserId(userId int) ([]*UserGroup, error)
}
type UserGroupRepo struct {
	DB *sql.DB
}

func NewUserGroupRepository(db *sql.DB) *UserGroupRepo {
	return &UserGroupRepo{
		DB: db,
	}
}

func (m *UserGroupRepo) GetAllUsergroups() ([]*UserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetAllUsergroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usergroups []*UserGroup
	for rows.Next() {
		usergroup, err := scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}

		usergroups = append(usergroups, usergroup)
	}
	return usergroups, nil
}
func (m *UserGroupRepo) GetById(id int) (*UserGroup, error) {
	rows, err := m.DB.Query(sqlGetById, id)
	if err != nil {
		return nil, err
	}

	usergroup := new(UserGroup)
	for rows.Next() {
		usergroup, err = scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}
	}
	return usergroup, err
}
func (m *UserGroupRepo) GetUserGroupsByUserId(userId int) ([]*UserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetUserGroupsByUserId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usergroups []*UserGroup
	for rows.Next() {
		usergroup, err := scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}

		usergroups = append(usergroups, usergroup)
	}
	return usergroups, nil
}

func scanRowsIntoUsergroup(rows *sql.Rows) (*UserGroup, error) {
	usergroup := new(UserGroup)
	err := rows.Scan(
		&usergroup.GroupId,
		&usergroup.GroupName,
		&usergroup.GroupDesc,
		&usergroup.Type,
		&usergroup.Thumbnail,
		&usergroup.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return usergroup, nil
}
