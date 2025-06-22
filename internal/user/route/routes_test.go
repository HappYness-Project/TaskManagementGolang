package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/mocks"
	userModel "github.com/happYness-Project/taskManagementGolang/internal/user/model"
	userGroupModel "github.com/happYness-Project/taskManagementGolang/internal/usergroup/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func TestUserHandler(t *testing.T) {

	env := configs.Env{}
	logger := loggers.Setup(env)
	mockUserRepo := new(mocks.MockUserRepo)
	mockUserGroupRepo := new(mocks.MockUserGroupRepo)
	handler := NewHandler(logger, mockUserRepo, mockUserGroupRepo)

	// Reset mocks before each test to prevent interference
	t.Cleanup(func() {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
	})

	t.Run("when get all users, Then return status code 200 with users", func(t *testing.T) {
		expectedUsers := []*userModel.User{{Id: 1, UserId: "1", UserName: "Alice"}}

		mockUserRepo.On("GetAllUsers").Return(expectedUsers, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/users/", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()
		router.Get("/api/users/", handler.handleGetUsers)
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify response structure
		assert.Equal(t, "success", response["message"])
		assert.NotNil(t, response["data"])

		// Parse data as array of users
		dataBytes, _ := json.Marshal(response["data"])
		var users []userModel.User
		err = json.Unmarshal(dataBytes, &users)
		require.NoError(t, err)

		// Verify user data
		assert.Len(t, users, 1)
		assert.Equal(t, expectedUsers[0].Id, users[0].Id)
		assert.Equal(t, expectedUsers[0].UserId, users[0].UserId)
		assert.Equal(t, expectedUsers[0].UserName, users[0].UserName)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("when get user by ID, Then return status code 200 with user details", func(t *testing.T) {
		userID := "test-user-id"
		expectedUser := &userModel.User{
			Id:        1,
			UserId:    userID,
			UserName:  "testuser",
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		expectedUserGroups := []*userGroupModel.UserGroup{
			{GroupId: 1, GroupName: "Admin", GroupDesc: "Administrators"},
		}

		mockUserRepo.On("GetUserByUserId", userID).Return(expectedUser, nil)
		mockUserGroupRepo.On("GetUserGroupsByUserId", expectedUser.Id).Return(expectedUserGroups, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		// Set up chi router with URL parameter
		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		dataBytes, _ := json.Marshal(response["data"])
		var userDetail UserDetailDto
		_ = json.Unmarshal(dataBytes, &userDetail)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, expectedUser.Id, userDetail.Id)
		assert.Equal(t, expectedUser.UserId, userDetail.UserId)
		assert.Equal(t, expectedUser.UserName, userDetail.UserName)
		assert.Equal(t, expectedUser.FirstName, userDetail.FirstName)
		assert.Equal(t, expectedUser.LastName, userDetail.LastName)
		assert.Equal(t, expectedUser.Email, userDetail.Email)
		assert.Equal(t, expectedUser.IsActive, userDetail.IsActive)
		assert.Equal(t, expectedUser.DefaultGroupId, userDetail.DefaultGroupId)

		assert.Len(t, userDetail.UserGroup, 1)
		assert.Equal(t, expectedUserGroups[0].GroupId, userDetail.UserGroup[0].GroupId)
		assert.Equal(t, expectedUserGroups[0].GroupName, userDetail.UserGroup[0].GroupName)
		assert.Equal(t, expectedUserGroups[0].GroupDesc, userDetail.UserGroup[0].GroupDesc)
	})

	t.Run("when get user by ID but error from repository layer, Then return status code 500", func(t *testing.T) {
		userID := "non-existent-user"

		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), fmt.Errorf("database error"))

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "database error")
	})

	t.Run("when get user by ID but user not found, Then return status code 404", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserGroupRepo.ExpectedCalls = nil
		userID := "non-existent-user"
		mockUserRepo.On("GetUserByUserId", userID).Return((*userModel.User)(nil), nil)

		req := httptest.NewRequest(http.MethodGet, "/api/users/"+userID, nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/users/{userID}", func(r chi.Router) {
			r.Get("/", handler.handleGetUser)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Parse error response
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "user does not exist")
	})

	t.Run("when get users by group ID, Then return status code 200 with users", func(t *testing.T) {
		groupID := "1"
		expectedUsers := []*userModel.User{
			{Id: 1, UserId: "user1", UserName: "Alice"},
			{Id: 2, UserId: "user2", UserName: "Bob"},
		}

		mockUserRepo.On("GetUsersByGroupId", 1).Return(expectedUsers, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/user-groups/"+groupID+"/users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Parse response body
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		dataBytes, _ := json.Marshal(response["data"])
		var users []userModel.User
		err = json.Unmarshal(dataBytes, &users)
		require.NoError(t, err)

		assert.Len(t, users, 2)
		assert.Equal(t, expectedUsers[0].Id, users[0].Id)
		assert.Equal(t, expectedUsers[0].UserId, users[0].UserId)
		assert.Equal(t, expectedUsers[0].UserName, users[0].UserName)
		assert.Equal(t, expectedUsers[1].Id, users[1].Id)
		assert.Equal(t, expectedUsers[1].UserId, users[1].UserId)
		assert.Equal(t, expectedUsers[1].UserName, users[1].UserName)
	})

	t.Run("when get users by group ID with invalid group ID, Then return status code 400", func(t *testing.T) {
		groupID := "invalid"

		req := httptest.NewRequest(http.MethodGet, "/api/user-groups/"+groupID+"/users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Parse error response
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "invalid user ID")
	})

	t.Run("when get users by group ID with empty group ID, Then return status code 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user-groups//users", nil)
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.Route("/api/user-groups/{groupID}", func(r chi.Router) {
			r.Get("/users", handler.handleGetUsersByGroupId)
		})

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Parse error response
		var response map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["error"].(bool))
		assert.Contains(t, response["message"], "missing Group ID")
	})
}
