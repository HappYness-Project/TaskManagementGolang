package mocks

import (
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/model"
	"github.com/stretchr/testify/mock"
)

type MockContainerRepo struct{ mock.Mock }

// AllTaskContainers implements repository.ContainerRepository.
func (m *MockContainerRepo) AllTaskContainers() ([]*model.TaskContainer, error) {
	args := m.Called()
	return args.Get(0).([]*model.TaskContainer), args.Error(1)
}

// GetById implements repository.ContainerRepository.
func (m *MockContainerRepo) GetById(id string) (*model.TaskContainer, error) {
	args := m.Called(id)
	return args.Get(0).(*model.TaskContainer), args.Error(1)
}

// GetContainersByGroupId implements repository.ContainerRepository.
func (m *MockContainerRepo) GetContainersByGroupId(groupId int) ([]model.TaskContainer, error) {
	args := m.Called(groupId)
	return args.Get(0).([]model.TaskContainer), args.Error(1)
}

// CreateContainer implements repository.ContainerRepository.
func (m *MockContainerRepo) CreateContainer(container model.TaskContainer) error {
	args := m.Called(container)
	return args.Error(0)
}

// DeleteContainer implements repository.ContainerRepository.
func (m *MockContainerRepo) DeleteContainer(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// RemoveContainerByUsergroupId implements repository.ContainerRepository.
func (m *MockContainerRepo) RemoveContainerByUsergroupId(groupId int) error {
	args := m.Called(groupId)
	return args.Error(0)
}
