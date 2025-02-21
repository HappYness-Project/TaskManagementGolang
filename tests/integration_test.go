package tests

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/happYness-Project/taskManagementGolang/cmd/api"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/dbs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func runTestServer() *httptest.Server {
	env := configs.InitConfig("integration_test.env")
	configs.AccessToken = env.AccessTokenSecret
	var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		env.DBHost, env.DBPort, env.DBUser, env.DBPwd, env.DBName)
	logger := loggers.Setup(env)
	logger.Info().Msg(connStr)
	database, err := dbs.ConnectToDb(connStr)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(fmt.Sprintf(":%d", 8000), database, logger)
	return httptest.NewServer(server.Setup())
}

func TestIntegrationTestHomePage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ts := runTestServer()
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", "http://localhost:8080"))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected 200 got: %v", resp.StatusCode)
	}
}
