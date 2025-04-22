package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
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

	// Create the mock logger
	logger := zaptest.NewLogger(t)

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

	// Create the mock logger
	logger := zaptest.NewLogger(t)

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

	mockFooRepo.EXPECT().
		CreateFoo().
		Return(int64(1), nil)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

	rowsAffected, err := fooService.CreateFoo()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)
}

func TestFooService_CreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooRepoError := errors.New("insert failed")
	mockFooRepo.EXPECT().
		CreateFoo().
		Return(int64(0), fooRepoError)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

	rowsAffected, err := fooService.CreateFoo()
	require.Equal(t, int64(0), rowsAffected)
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

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

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

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

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
	mockFooRepo.EXPECT().
		UpdateFoo(fooID).
		Return(int64(1), nil)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

	rowsAffected, err := fooService.UpdateFoo(fooID)
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)
}

func TestFooService_UpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFooRepo := mocks.NewMockFooRepo(ctrl)

	fooID := int64(100)
	fooRepoError := errors.New("update failed")
	mockFooRepo.EXPECT().
		UpdateFoo(fooID).
		Return(int64(0), fooRepoError)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	fooService := NewFooService(mockFooRepo, logger)

	rowsAffected, err := fooService.UpdateFoo(fooID)
	require.Equal(t, int64(0), rowsAffected)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GZNHKW")
}
