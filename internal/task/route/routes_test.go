package route

// func TestTaskHandler(t *testing.T) {
// 	mockTaskRepo := new(mocks.TaskRepository)
// 	mockContainerRepo := new(mocks.ContainerRepository)

// 	env := configs.Env{}
// 	logger := loggers.Setup(env)

// 	handler := NewHandler(logger, mockTaskRepo, mockContainerRepo, nil)

// 	t.Run("when get all tasks, Then return status code 200 with tasks", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Get("/api/tasks", handler.handleGetTasks)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusOK)
// 		var tasks []model.Task
// 		err := json.Unmarshal(rr.Body.Bytes(), &tasks)
// 		require.NoError(t, err)
// 		assert.NotNil(t, tasks)
// 	})
// }

// 	t.Run("when get task by id, Then return status code 200", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodGet, "/api/tasks/abcd", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Get("/api/tasks/{taskID}", handler.handleGetTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusOK)
// 	})
// 	t.Run("given missing container Id, when handleGetTasksByContainerId called, Then bad request", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodGet, "/api/task-containers//tasks", nil)
// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Get("/api/task-containers/{containerID}/tasks", handler.handleGetTasksByContainerId)
// 		mux.ServeHTTP(rr, req)

// 		assertStatus(t, rr.Code, http.StatusBadRequest)
// 	})
// 	t.Run("when create new task, then return status code 201", func(t *testing.T) {
// 		newTask := CreateTaskDto{
// 			TaskName:   "New task",
// 			TaskDesc:   "desc",
// 			TargetDate: time.Now().AddDate(0, 0, 7*1),
// 			Priority:   "normal",
// 			Category:   "programming",
// 		}
// 		marshalled, _ := json.Marshal(newTask)
// 		req, _ := http.NewRequest(http.MethodPost, "/api/task-containers/5951f639-c8ce-4462-8b72-c57458c448fd/tasks", bytes.NewBuffer(marshalled))
// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Post("/api/task-containers/{containerID}/tasks", handler.handleCreateTask)
// 		mux.ServeHTTP(rr, req)

// 		assertStatus(t, rr.Code, http.StatusCreated)
// 	})
// 	t.Run("given payload is missing, when creating new task, then bad request", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodPost, "/api/task-containers/5951f639-c8ce-4462-8b72-c57458c448fd/tasks", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Post("/api/task-containers/{containerID}/tasks", handler.handleCreateTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusBadRequest)
// 	})
// 	t.Run("given valid payload, when updating existing task, then success", func(t *testing.T) {
// 		updateTask := UpdateTaskDto{
// 			TaskName:   "New task",
// 			TaskDesc:   "desc",
// 			TargetDate: time.Now().AddDate(0, 0, 7*1),
// 			Priority:   "normal",
// 			Category:   "programming",
// 		}
// 		marshalled, _ := json.Marshal(updateTask)
// 		req, _ := http.NewRequest(http.MethodPut, "/api/tasks/5951f639-c8ce-4462-8b72-c57458c448fd", bytes.NewBuffer(marshalled))

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Put("/api/tasks/{taskID}", handler.handleUpdateTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusOK)
// 	})
// 	t.Run("given payload is missing, when updating existing task, then bad request", func(t *testing.T) {
// 		req, err := http.NewRequest(http.MethodPut, "/api/tasks/5951f639-c8ce-4462-8b72-c57458c448fd", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Put("/api/tasks/{taskID}", handler.handleUpdateTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusBadRequest)
// 	})
// 	t.Run("given empty ID, when delete task, then not found", func(t *testing.T) {
// 		req, _ := http.NewRequest(http.MethodDelete, "/api/tasks/", nil)

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Delete("/api/tasks/{taskID}", handler.handleDeleteTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusNotFound)
// 	})
// 	t.Run("given empty taskID format, when togglecompletion, then bad request", func(t *testing.T) {
// 		var jsonStr = []byte(`{"is_completed":true }`)
// 		req, _ := http.NewRequest(http.MethodPatch, "/api/tasks//toggle-completion", bytes.NewBuffer(jsonStr))

// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Patch("/api/tasks/{taskID}/toggle-completion", handler.handleDoneTask)
// 		mux.ServeHTTP(rr, req)
// 		assertStatus(t, rr.Code, http.StatusBadRequest)
// 	})
// 	t.Run("given identifier with body, when togglecompletion, then success", func(t *testing.T) {
// 		var jsonStr = []byte(`{"is_completed":true }`)
// 		req, _ := http.NewRequest(http.MethodPatch, "/api/tasks/5951f639-c8ce-4462-8b72-c57458c448fd/toggle-completion", bytes.NewBuffer(jsonStr))
// 		rr := httptest.NewRecorder()
// 		mux := chi.NewRouter()

// 		mux.Patch("/api/tasks/{taskID}/toggle-completion", handler.handleDoneTask)
// 		mux.ServeHTTP(rr, req)

// 		assertStatus(t, rr.Code, http.StatusOK)
// 	})
// }

// func (m *mockTaskRepo) GetAllTasks() ([]model.Task, error) {
// 	var mockTasks = []model.Task{
// 		{TaskId: uuid.New().String(), TaskName: "TaskName #1", TaskDesc: "Task Desc #1", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*1)},
// 		{TaskId: uuid.New().String(), TaskName: "TaskName #2", TaskDesc: "Task Desc #2", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*2)},
// 		{TaskId: uuid.New().String(), TaskName: "TaskName #3", TaskDesc: "Task Desc #3", TaskType: "", TargetDate: time.Now().AddDate(0, 0, 7*3)},
// 	}
// 	return mockTasks, nil
// }

// func assertStatus(t testing.TB, got, want int) {
// 	t.Helper()
// 	if got != want {
// 		t.Errorf("expected status code %d, got %d", got, want)
// 	}
// }
