package task

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer"
)

func TestTaskHandler(t *testing.T) {
	repo := &mockTaskRepo{}
	containerRepo := &mockContainerRepo{}
	handler := NewHandler(repo, containerRepo)

	t.Run("when get all tasks, Then return status code 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Get("/api/tasks", handler.handleGetTasks)

		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
	t.Run("when get task by id, Then return status code 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/tasks/abcd", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Get("/api/tasks/{taskID}", handler.handleGetTask)

		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockTaskRepo struct{}

func (m *mockTaskRepo) GetAllTasks() ([]*Task, error) {
	return []*Task{}, nil
}
func (m *mockTaskRepo) GetTaskById(id string) (*Task, error) {
	return &Task{}, nil
}
func (m *mockTaskRepo) GetTasksByContainerId(containerId string) ([]*Task, error) {
	return []*Task{}, nil
}
func (m *mockTaskRepo) CreateTask(containerId string, req CreateTaskDto) (int64, error) {
	return 0, nil
}
func (m *mockTaskRepo) UpdateTask(task UpdateTaskDto) error {
	return nil
}
func (m *mockTaskRepo) DeleteTask(id string) error {
	return nil
}
func (m *mockTaskRepo) DoneTask(id string, isDone bool) error {
	return nil
}
func (m *mockTaskRepo) UpdateImportantTask(id string) error {
	return nil
}

type mockContainerRepo struct{}

func (m *mockContainerRepo) AllTaskContainers() ([]*taskcontainer.TaskContainer, error) {
	return []*taskcontainer.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetById(id string) (*taskcontainer.TaskContainer, error) {
	return &taskcontainer.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetContainersByGroupId(groupId int) ([]*taskcontainer.TaskContainer, error) {
	return []*taskcontainer.TaskContainer{}, nil
}
