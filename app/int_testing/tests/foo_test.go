package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func TestFooRepo_GetFoos_Success(t *testing.T) {
	ctx := GetContext()

	// Get app from context
	appVal := (*ctx).Value("App")
	app, ok := appVal.(*fiber.App)
	require.True(t, ok, "App not found in context or wrong type")
	require.NotNil(t, app, "App is nil")

	// Build request
	req := httptest.NewRequest(http.MethodGet, "/foos", nil)

	// Run request
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var foos []models.Foo
	err = json.Unmarshal(body, &foos)
	require.NoError(t, err)

	require.Len(t, foos, 3)
	require.Equal(t, "Test Foo 1", foos[0].Name)
	require.Equal(t, 2, foos[1].ID)
}
