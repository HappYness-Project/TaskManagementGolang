package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/mocks"
	userModel "github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/stretchr/testify/assert"

	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

// func (m *mock.MockUserRepo) GetAllUsers() ([]*userModel.User, error) {
// 	args := m.Called()
// 	return args.Get(0).([]*userModel.User), args.Error(1)
// }

// func (m *MockUserGroupRepo) GetUserGroupByID(ctx context.Context, id int) (*userGroupModel.UserGroup, error) {
// 	args := m.Called(ctx, id)
// 	return args.Get(0).(*userGroupModel.UserGroup), args.Error(1)
// }

func TestUserHandler(t *testing.T) {

	env := configs.Env{}
	logger := loggers.Setup(env)
	mockUserRepo := new(mocks.MockUserRepo)
	mockUserGroupRepo := new(mocks.MockUserGroupRepo)
	handler := NewHandler(logger, mockUserRepo, mockUserGroupRepo)

	expectedUsers := []*userModel.User{{Id: 1, UserId: "1", UserName: "Alice"}}

	mockUserRepo.On("GetAllUsers").Return(expectedUsers, nil)
	req := httptest.NewRequest(http.MethodGet, "/api/users/", nil)
	rr := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Get("/api/users/", handler.handleGetUsers)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
