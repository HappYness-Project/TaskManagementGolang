package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler(t *testing.T) {
	repo := &mockTaskRepo{}
	containerRepo := &mockContainerRepo{}
	handler := NewHandler(repo, containerRepo, nil)

	t.Run("when get all tasks, Then return status code 200 with tasks", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Get("/api/tasks", handler.handleGetTasks)
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var tasks []Task
		err := json.Unmarshal(rr.Body.Bytes(), &tasks)
		require.NoError(t, err)
		assert.NotNil(t, tasks)
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
	t.Run("when create new task, then return status code 201", func(t *testing.T) {
		newTask := CreateTaskDto{
			TaskName:   "New task",
			TaskDesc:   "desc",
			TargetDate: time.Now().AddDate(0, 0, 7*1),
			Priority:   "normal",
			Category:   "programming",
		}
		marshalled, _ := json.Marshal(newTask)
		req, err := http.NewRequest(http.MethodPost, "/api/task-containers/5951f639-c8ce-4462-8b72-c57458c448fd/tasks", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Post("/api/task-containers/{containerID}/tasks", handler.handleCreateTask)
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
	t.Run("when payload is missing, then bad request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/task-containers/5951f639-c8ce-4462-8b72-c57458c448fd/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Post("/api/task-containers/{containerID}/tasks", handler.handleCreateTask)
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockTaskRepo struct{}

func (m *mockTaskRepo) GetAllTasks() ([]Task, error) {
	var mockTasks = []Task{
		{TaskId: uuid.New().String(), TaskName: "TaskName #1", TaskDesc: "Task Desc #1", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*1)},
		{TaskId: uuid.New().String(), TaskName: "TaskName #2", TaskDesc: "Task Desc #2", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*2)},
		{TaskId: uuid.New().String(), TaskName: "TaskName #3", TaskDesc: "Task Desc #3", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*3)},
	}
	return mockTasks, nil
}
func (m *mockTaskRepo) GetTaskById(id string) (*Task, error) {
	return &Task{}, nil
}
func (m *mockTaskRepo) GetTasksByContainerId(containerId string) ([]Task, error) {
	return []Task{}, nil
}
func (m *mockTaskRepo) CreateTask(containerId string, req Task) (Task, error) {
	return Task{}, nil
}
func (m *mockTaskRepo) UpdateTask(task Task) error {
	return nil
}
func (m *mockTaskRepo) DeleteTask(id string) error {
	return nil
}
func (m *mockTaskRepo) DoneTask(id string, isDone bool) error {
	return nil
}
func (m *mockTaskRepo) UpdateImportantTask(id string, isImportant bool) error {
	return nil
}
func (m *mockTaskRepo) GetAllTasksByGroupId(groupId int) ([]Task, error) {
	return []Task{}, nil
}
func (m *mockTaskRepo) GetAllTasksByGroupIdOnlyImportant(groupId int) ([]Task, error) {
	return []Task{}, nil
}

type mockContainerRepo struct{}

func (m *mockContainerRepo) AllTaskContainers() ([]*taskcontainer.TaskContainer, error) {
	return []*taskcontainer.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetById(id string) (*taskcontainer.TaskContainer, error) {
	return &taskcontainer.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetContainersByGroupId(groupId int) ([]taskcontainer.TaskContainer, error) {
	return []taskcontainer.TaskContainer{}, nil
}
func (m *mockContainerRepo) CreateContainer(container taskcontainer.TaskContainer) error {
	return nil
}
