package mocks

import (
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/model"
	"github.com/stretchr/testify/mock"
)

type MockUserGroupRepo struct{ mock.Mock }

// CreateGroup implements repository.UserGroupRepository.
func (m MockUserGroupRepo) CreateGroup(ug model.UserGroup) (int, error) {
	panic("unimplemented")
}

// DeleteUserGroup implements repository.UserGroupRepository.
func (m MockUserGroupRepo) DeleteUserGroup(id int) error {
	panic("unimplemented")
}

// GetAllUsergroups implements repository.UserGroupRepository.
func (m MockUserGroupRepo) GetAllUsergroups() ([]*model.UserGroup, error) {
	panic("unimplemented")
}

// GetById implements repository.UserGroupRepository.
func (m MockUserGroupRepo) GetById(id int) (*model.UserGroup, error) {
	panic("unimplemented")
}

// GetUserGroupsByUserId implements repository.UserGroupRepository.
func (m MockUserGroupRepo) GetUserGroupsByUserId(userId int) ([]*model.UserGroup, error) {
	panic("unimplemented")
}

// InsertUserGroupUserTable implements repository.UserGroupRepository.
func (m MockUserGroupRepo) InsertUserGroupUserTable(groupId int, userId int) error {
	panic("unimplemented")
}

// RemoveUserFromUserGroup implements repository.UserGroupRepository.
func (m MockUserGroupRepo) RemoveUserFromUserGroup(groupId int, userId int) error {
	panic("unimplemented")
}
