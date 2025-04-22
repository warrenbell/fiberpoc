// common/handler/foo_handler_test.go
package handlers

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/sandstone2/fiberpoc/common/mocks"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func TestFooHandler_HandleGetFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. Create mock service and logger
	mockFooService := mocks.NewMockFooService(ctrl)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	// 2. Instantiate handler
	fooHandler := NewFooHandler(mockFooService, logger)

	// 3. Set up a Fiber app and route
	app := fiber.New()
	app.Get("/foos", fooHandler.HandleGetFoos)

	// 4. Stub service to return a known value
	expected := &[]models.Foo{
		{ID: 1, Name: "Foo One"},
		{ID: 2, Name: "Foo Two"},
	}
	mockFooService.
		EXPECT().
		GetFoos().
		Return(expected, nil)

	// 5. Perform the HTTP request
	request := httptest.NewRequest("GET", "/foos", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 6. Assertions
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	// JSONEq ignores key order, spacing, etc.
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
	require.JSONEq(t, `{"message":"Error J5TSGF - Getting foos in handler."}`, string(body))
}

func TestFooHandler_HandleCreateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. Mock the service and create handler
	mockFooService := mocks.NewMockFooService(ctrl)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooHandler := NewFooHandler(mockFooService, logger)

	// 2. Set up Fiber and route
	app := fiber.New()
	app.Post("/foo", fooHandler.HandleCreateFoo)

	// 3. Expect CreateFoo to be called, returning 3 rows inserted
	mockFooService.EXPECT().
		CreateFoo().
		Return(int64(1), nil)

	// 4. Perform the request
	request := httptest.NewRequest("POST", "/foo", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assert 200 OK and the correct JSON message
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	expected := `{"message":"1 foos created."}`
	require.JSONEq(t, expected, string(body))
}

func TestFooHandler_HandleCreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. Mock the service and create handler
	mockFooService := mocks.NewMockFooService(ctrl)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooHandler := NewFooHandler(mockFooService, logger)

	// 2. Set up Fiber and route
	app := fiber.New()
	app.Post("/foo", fooHandler.HandleCreateFoo)

	// 3. Expect CreateFoo to return an error
	mockFooService.EXPECT().
		CreateFoo().
		Return(int64(0), errors.New("fail"))

	// 4. Perform the request
	request := httptest.NewRequest("POST", "/foo", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assert 500 and the predefined error JSON
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	require.JSONEq(t,
		`{"message":"Error QONMRA - Creating foo in handler."}`,
		string(body),
	)
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

	// 1. Mock the service and create handler
	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	// 2. Set up Fiber and route
	app := fiber.New()
	app.Patch("/foo", fooHandler.HandleUpdateFoo)

	// 3. Expect UpdateFoo to be called with fooId=42, returning 2 rows updated
	mockFooService.
		EXPECT().
		UpdateFoo(int64(42)).
		Return(int64(1), nil)

	// 4. Perform the request with ?fooId=42
	request := httptest.NewRequest("PATCH", "/foo?fooId=42", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assert 200 OK and the correct JSON message
	require.Equal(t, fiber.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	expected := `{"message":"1 foos updated."}`
	require.JSONEq(t, expected, string(body))
}

func TestFooHandler_HandleUpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooService := mocks.NewMockFooService(ctrl)
	logger := zaptest.NewLogger(t)
	fooHandler := NewFooHandler(mockFooService, logger)

	app := fiber.New()
	app.Patch("/foo", fooHandler.HandleUpdateFoo)

	// 3. Stub UpdateFoo to return an error
	mockFooService.
		EXPECT().
		UpdateFoo(int64(42)).
		Return(int64(0), errors.New("fail"))

	// 4. Perform the request
	request := httptest.NewRequest("PATCH", "/foo?fooId=42", nil)
	response, err := app.Test(request, -1)
	require.NoError(t, err)
	defer response.Body.Close()

	// 5. Assert 500 and the predefined error JSON
	require.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
	body, _ := io.ReadAll(response.Body)
	require.JSONEq(t,
		`{"message":"Error E4LP9X - Updating foo in handler."}`,
		string(body),
	)
}
