package route

import (
	"net/http"
	"net/http/httptest"
	"os/user"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/model"
)

func TestTaskContainerHandler(t *testing.T) {
	repo := &mockContainerRepo{}
	// userRepo := &mockUserRepo{}
	handler := NewHandler(repo, nil) // Need to be fixed.

	t.Run("when get all task containers, Then return status code 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/task-containers", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Get("/api/task-containers", handler.handleGetTaskContainers)

		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
	t.Run("when get task container by containerId, Then return status code 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/task-containers/abcd", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		mux := chi.NewRouter()

		mux.Get("/api/task-containers/{containerID}", handler.handleGetTaskContainerById)

		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		defer rr.Result().Body.Close()
	})
}

type mockContainerRepo struct{}

func (m *mockContainerRepo) AllTaskContainers() ([]*model.TaskContainer, error) {
	return []*model.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetById(id string) (*model.TaskContainer, error) {
	return &model.TaskContainer{}, nil
}
func (m *mockContainerRepo) GetContainersByGroupId(groupId int) ([]model.TaskContainer, error) {
	return []model.TaskContainer{}, nil
}
func (m *mockContainerRepo) CreateContainer(c model.TaskContainer) error {
	return nil
}
func (m *mockContainerRepo) DeleteContainer(id string) error {
	return nil
}

type mockUserRepo struct{}

func (m *mockUserRepo) GetAllUsers() ([]*user.User, error) {
	return []*user.User{}, nil
}
