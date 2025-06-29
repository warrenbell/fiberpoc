package services

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"gitlab.com/sandstone2/fiberpoc/common/mocks"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func TestFooService_GetFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	expected := &[]models.Foo{{ID: 1, Name: "Joe"}}
	mockFooRepo.EXPECT().
		GetFoos().
		Return(expected, nil)

	logger := zaptest.NewLogger(t)

	// fix: pass a pointer to mockFooRepo
	fooService := NewFooService(mockFooRepo, logger)

	foos, err := fooService.GetFoos()
	require.NoError(t, err)
	require.Equal(t, expected, foos)
}

func TestFooService_GetFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooRepoError := errors.New("db failure")
	mockFooRepo.EXPECT().
		GetFoos().
		Return(nil, fooRepoError)

	logger := zaptest.NewLogger(t)

	// Pass pointer to mockFooRepo
	fooService := NewFooService(mockFooRepo, logger)

	foos, err := fooService.GetFoos()
	require.Nil(t, foos)
	require.Error(t, err)
	require.Contains(t, err.Error(), "WZDCXT")
}

func TestFooService_CreateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	// Set up expected Foo to return
	expectedFoo := &models.Foo{ID: 1, Name: "Test Foo"}

	mockFooRepo.EXPECT().
		CreateFoo("Test Foo").
		Return(expectedFoo, nil)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // Pass pointer

	foo, err := fooService.CreateFoo("Test Foo")
	require.NoError(t, err)
	require.Equal(t, expectedFoo, foo)
}

func TestFooService_CreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooRepoError := errors.New("insert failed")
	mockFooRepo.EXPECT().
		CreateFoo("Test Foo").
		Return(nil, fooRepoError)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // pass pointer

	foo, err := fooService.CreateFoo("Test Foo")
	require.Nil(t, foo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "DWA4G7")
}

func TestFooService_DeleteFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	mockFooRepo.EXPECT().
		DeleteFoos().
		Return(int64(5), nil)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // pass pointer

	rowsAffected, err := fooService.DeleteFoos()
	require.NoError(t, err)
	require.Equal(t, int64(5), rowsAffected)
}

func TestFooService_DeleteFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooRepoError := errors.New("delete failed")
	mockFooRepo.EXPECT().
		DeleteFoos().
		Return(int64(0), fooRepoError)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // pass pointer

	rowsAffected, err := fooService.DeleteFoos()
	require.Equal(t, int64(0), rowsAffected)
	require.Error(t, err)
	require.Contains(t, err.Error(), "BA8TAX")
}

func TestFooService_UpdateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooID := int64(42)
	newName := "Updated Name"

	expectedFoo := &models.Foo{ID: int(fooID), Name: newName}

	mockFooRepo.EXPECT().
		UpdateFoo(fooID, newName).
		Return(expectedFoo, nil)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // pass pointer

	foo, err := fooService.UpdateFoo(fooID, newName)
	require.NoError(t, err)
	require.Equal(t, expectedFoo, foo)
}

func TestFooService_UpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooID := int64(100)
	newName := "Some Name"
	fooRepoError := errors.New("update failed")

	mockFooRepo.EXPECT().
		UpdateFoo(fooID, newName).
		Return(nil, fooRepoError)

	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger) // pass pointer

	foo, err := fooService.UpdateFoo(fooID, newName)
	require.Nil(t, foo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GZNHKW")
}
