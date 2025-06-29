// common/handler/foo_handler_test.go
package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"gitlab.com/sandstone2/fiberpoc/common/mocks"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func TestFooHandler_HandleGetFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)

	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Get("/foos", fooHandler.HandleGetFoos)

	expected := &[]models.Foo{
		{ID: 1, Name: "Foo One"},
		{ID: 2, Name: "Foo Two"},
	}
	mockFooService.
		EXPECT().
		GetFoos().
		Return(expected, nil)

	request := httptest.NewRequest("GET", "/foos", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	require.JSONEq(t, `[{"ID":1,"Name":"Foo One"},{"ID":2,"Name":"Foo Two"}]`, string(body))
}

func TestFooHandler_HandleGetFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Get("/foos", fooHandler.HandleGetFoos)

	// Stub service to return an error
	serviceErr := errors.New("db failure")
	mockFooService.
		EXPECT().
		GetFoos().
		Return(nil, serviceErr)

	request := httptest.NewRequest("GET", "/foos", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	require.JSONEq(t, `{"message":"Error J5TSGF - Getting foos in handler. Error: db failure"}`, string(body))
}

func TestFooHandler_HandleCreateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Post("/foo", fooHandler.HandleCreateFoo)

	// Prepare the input Foo JSON
	inputJSON := `{"Name":"New Foo"}`
	request := httptest.NewRequest("POST", "/foo", strings.NewReader(inputJSON))
	request.Header.Set("Content-Type", "application/json")

	// Expected Foo to be returned from the service
	createdFoo := &models.Foo{ID: 1, Name: "New Foo"}

	// Expect CreateFoo(name) to be called with "New Foo" and return createdFoo
	mockFooService.
		EXPECT().
		CreateFoo("New Foo").
		Return(createdFoo, nil)

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	// The handler returns the created Foo object as JSON
	expectedJSON := `{"ID":1,"Name":"New Foo"}`
	require.JSONEq(t, expectedJSON, string(body))
}

func TestFooHandler_HandleCreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Post("/foo", fooHandler.HandleCreateFoo)

	inputJSON := `{"Name":"Bad Foo"}`
	request := httptest.NewRequest("POST", "/foo", strings.NewReader(inputJSON))
	request.Header.Set("Content-Type", "application/json")

	expectedErr := errors.New("fail")
	mockFooService.
		EXPECT().
		CreateFoo("Bad Foo").
		Return(nil, expectedErr)

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	// Your handler formats error with the err string
	expectedMessage := fmt.Sprintf(`{"message":"Error QONMRA - Creating foo in handler. Error: %v"}`, expectedErr)
	require.JSONEq(t, expectedMessage, string(body))
}

func TestFooHandler_HandleDeleteFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. Create mock service and logger
	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	// 2. Set up a Fiber app and route
	app := fiber.New()
	app.Delete("/foos", fooHandler.HandleDeleteFoos)

	// 3. Stub service to return 5 rows deleted
	mockFooService.
		EXPECT().
		DeleteFoos().
		Return(int64(5), nil)

	// 4. Perform the HTTP request
	request := httptest.NewRequest("DELETE", "/foos", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assertions
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	expected := `{"message":"5 foos deleted."}`
	require.JSONEq(t, expected, string(body))
}

func TestFooHandler_HandleDeleteFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Delete("/foo", fooHandler.HandleDeleteFoos)

	// 3. Stub service to return an error
	mockFooService.
		EXPECT().
		DeleteFoos().
		Return(int64(0), errors.New("fail"))

	// 4. Perform the HTTP request
	request := httptest.NewRequest("DELETE", "/foo", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assertions
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	require.JSONEq(t,
		`{"message":"Error 8HCIPG - Deleting foos in handler."}`,
		string(body),
	)
}

func TestFooHandler_HandleUpdateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	// Use route with :id param to match handler expectations
	app.Patch("/foo/:id", fooHandler.HandleUpdateFoo)

	// Prepare request body JSON with updated Foo name
	inputJSON := `{"Name":"Updated Foo"}`

	// Expect UpdateFoo to be called with id=42 and name="Updated Foo"
	expectedFoo := &models.Foo{ID: 42, Name: "Updated Foo"}

	mockFooService.
		EXPECT().
		UpdateFoo(int64(42), "Updated Foo").
		Return(expectedFoo, nil)

	request := httptest.NewRequest("PATCH", "/foo/42", strings.NewReader(inputJSON))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	expectedJSON := `{"ID":42,"Name":"Updated Foo"}`
	require.JSONEq(t, expectedJSON, string(body))
}

func TestFooHandler_HandleUpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Patch("/foo/:id", fooHandler.HandleUpdateFoo)

	inputJSON := `{"Name":"Updated Foo"}`
	request := httptest.NewRequest("PATCH", "/foo/42", strings.NewReader(inputJSON))
	request.Header.Set("Content-Type", "application/json")

	expectedErr := errors.New("fail")

	mockFooService.
		EXPECT().
		UpdateFoo(int64(42), "Updated Foo").
		Return(nil, expectedErr)

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	expectedMessage := fmt.Sprintf(`{"message":"Error FSYTGZ - Updating foo. Error: %v"}`, expectedErr)
	require.JSONEq(t, expectedMessage, string(body))
}
