package utils_test

import (
	"testing"
	"time"

	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

func Test_FormatTimeToISO_Success(t *testing.T) {
	got := utils.FormatTimeToISO(time.Date(2025, 1, 1, 9, 30, 0, 0, time.UTC))
	want := "2025-01-01T09:30:00Z"

	if want != got {
		t.Errorf("Expected '%s', but got '%s'", want, got)
	}
}
