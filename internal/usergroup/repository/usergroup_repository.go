package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/model"
)

const dbTimeout = time.Second * 5

type UserGroupRepository interface {
	GetAllUsergroups() ([]*model.UserGroup, error)
	GetById(id int) (*model.UserGroup, error)
	GetUserGroupsByUserId(userId int) ([]*model.UserGroup, error)
	CreateGroup(ug model.UserGroup) (int, error)
	InsertUserGroupUserTable(groupId int, userId int) error
	RemoveUserFromUserGroup(groupId int, userId int) error
	DeleteUserGroup(id int) error
}
type UserGroupRepo struct {
	DB *sql.DB
}

func NewUserGroupRepository(db *sql.DB) *UserGroupRepo {
	return &UserGroupRepo{
		DB: db,
	}
}

func (m *UserGroupRepo) GetAllUsergroups() ([]*model.UserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetAllUsergroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usergroups []*model.UserGroup
	for rows.Next() {
		usergroup, err := scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}

		usergroups = append(usergroups, usergroup)
	}
	return usergroups, nil
}
func (m *UserGroupRepo) GetById(id int) (*model.UserGroup, error) {
	rows, err := m.DB.Query(sqlGetById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usergroup := new(model.UserGroup)
	for rows.Next() {
		usergroup, err = scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}
	}
	return usergroup, err
}
func (m *UserGroupRepo) GetUserGroupsByUserId(userIntId int) ([]*model.UserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, sqlGetUserGroupsByUserId, userIntId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usergroups []*model.UserGroup
	for rows.Next() {
		usergroup, err := scanRowsIntoUsergroup(rows)
		if err != nil {
			return nil, err
		}

		usergroups = append(usergroups, usergroup)
	}
	return usergroups, nil
}
func (m *UserGroupRepo) CreateGroup(ug model.UserGroup) (int, error) {
	lastInsertedId := 0
	err := m.DB.QueryRow(sqlCreateUserGroup, ug.GroupName, ug.GroupDesc, ug.Type, ug.Thumbnail, ug.IsActive).Scan(&lastInsertedId)
	if err != nil {
		return 0, fmt.Errorf("unable to insert into usergroup table : %w", err)
	}

	return lastInsertedId, nil
}
func (m *UserGroupRepo) InsertUserGroupUserTable(groupId int, userId int) error {
	_, err := m.DB.Exec(sqlAddUserToUserGroup, groupId, userId)
	if err != nil {
		return fmt.Errorf("unable to insert into usergroup_user table : %w", err)
	}
	return nil
}

func (m *UserGroupRepo) RemoveUserFromUserGroup(groupId int, userId int) error {
	result, err := m.DB.Exec(sqlRemoveUserFromUserGroup, groupId, userId)
	if err != nil {
		return fmt.Errorf("unable to delete user from usergroup_user table : %w", err)
	}
	id, _ := result.RowsAffected()
	if id == 0 {
		return fmt.Errorf("none of the data has been removed")
	}
	return nil
}

func (m *UserGroupRepo) DeleteUserGroup(groupId int) error {
	_, err := m.DB.Exec(sqlDeleteUserGroup, groupId)
	if err != nil {
		return fmt.Errorf("unable to delete usergroup table : %w", err)
	}
	return nil
}

func scanRowsIntoUsergroup(rows *sql.Rows) (*model.UserGroup, error) {
	usergroup := new(model.UserGroup)
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
