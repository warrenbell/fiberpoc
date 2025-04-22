// myapp/repository/repository_test.go
package repos

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"gitlab.com/sandstone2/fiberpoc/common/mocks"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func TestFooRepo_GetFoos_Success(t *testing.T) {
	// Create a new GoMock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock for PgxPoolInterface.
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Create the mock for PgxRowsInterface.
	mockRows := mocks.NewMockPgxRows(ctrl)

	// Set expectation for the mock pgx pool.
	// Make sure the right query is called.
	const expectedQuery = "SELECT * FROM foos;"
	mockPool.EXPECT().
		Query(gomock.Any(), expectedQuery).
		Return(mockRows, nil)

	// Set expectations for the mock pgx rows.
	// We'll simulate that there is one row to return.
	// First, Next() returns true.
	mockRows.EXPECT().Next().Return(true)
	// When Scan() is called, we simulate scanning a row with ID = 1 and Name = "Foo One".
	mockRows.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			// Set dest[0] (pointer to int) to 1 and dest[1] (pointer to string) to "Foo One".
			*(dest[0].(*int)) = 1
			*(dest[1].(*string)) = "Joe"
			return nil
		})
	// After the one row, Next() returns false.
	mockRows.EXPECT().Next().Return(false)
	// rows.Err() returns nil.
	mockRows.EXPECT().Err().Return(nil)
	// rows.Close() is called.
	mockRows.EXPECT().Close()

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	// Create a ne foo repo
	fooRepo := NewFooRepository(mockPool, logger)

	// Call the GetFoos function under test.
	foos, err := fooRepo.GetFoos()
	require.NoError(t, err, "GetFoos should not return an error.")
	require.NotNil(t, foos, "foos should not be nil.")
	require.Len(t, *foos, 1, "foos should be length 1.")

	// Verify that the foos slice contains the expected result.
	expectedFoo := models.Foo{
		ID:   1,
		Name: "Joe",
	}
	require.Equal(t, expectedFoo, (*foos)[0], "foo returned should be correct.")
}

func TestFooRepo_GetFoos_Error(t *testing.T) {
	// There are three different error paths

	// Create a new GoMock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock for PgxPoolInterface.
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Create the mock for PgxRowsInterface.
	mockRows := mocks.NewMockPgxRows(ctrl)

	// 1) Test mockPool.Query failed
	// Set expectation for the mock pgx pool.
	// Make sure the right query is called.
	mockPool.EXPECT().
		Query(gomock.Any(), "SELECT * FROM foos;").
		Return(mockRows, errors.New("query failed"))

	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Call under test
	_, err := fooRepo.GetFoos()
	require.Error(t, err)
	require.Contains(t, err.Error(), "30UUBR", "error should be wrapped with 30UUBR code")

	// 2) Test mockRows.Scan failed
	mockPool.EXPECT().
		Query(gomock.Any(), "SELECT * FROM foos;").
		Return(mockRows, nil)

	mockRows.EXPECT().Next().Return(true)
	// When Scan() is called, we simulate scanning a row with ID = 1 and Name = "Foo One".
	mockRows.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		Return(errors.New("scan failed"))

	mockRows.EXPECT().Close()

	_, err = fooRepo.GetFoos()
	require.Error(t, err)
	require.Contains(t, err.Error(), "YN80XB", "error should be wrapped with YN80XB code")

	// 3) Test mockRows.Err failed
	mockPool.EXPECT().
		Query(gomock.Any(), "SELECT * FROM foos;").
		Return(mockRows, nil)

	// Set expectations for the mock pgx rows.
	// We'll simulate that there is one row to return.
	// First, Next() returns true.
	mockRows.EXPECT().Next().Return(true)
	// When Scan() is called, we simulate scanning a row with ID = 1 and Name = "Foo One".
	mockRows.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest ...interface{}) error {
			// Set dest[0] (pointer to int) to 1 and dest[1] (pointer to string) to "Foo One".
			*(dest[0].(*int)) = 1
			*(dest[1].(*string)) = "Joe"
			return nil
		})
	// After the one row, Next() returns false.
	mockRows.EXPECT().Next().Return(false)
	// rows.Err() returns nil.
	mockRows.EXPECT().Err().Return(errors.New("rows.Err failed"))
	// rows.Close() is called.
	mockRows.EXPECT().Close()

	_, err = fooRepo.GetFoos()
	require.Error(t, err)
	require.Contains(t, err.Error(), "XV4HHL", "error should be wrapped with XV4HHL code")
}

func TestFooRepo_CreateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Set expectation for the mock pgx pool.
	// Make sure the right insert is called.
	const insertSQL = "INSERT INTO foos (name) VALUES ($1);"
	mockPool.EXPECT().
		Exec(
			gomock.Any(),
			insertSQL,
			gomock.Any(),
		).
		Return(pgconn.NewCommandTag("INSERT 1"), nil)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	// Create a ne foo repo
	fooRepo := NewFooRepository(mockPool, logger)

	// Call the CreateFoo function under test.
	rows, err := fooRepo.CreateFoo()
	require.NoError(t, err, "CreateFoo should not return an error.")
	require.Equal(t, int64(1), rows, "should report 1 row affected.")
}

func TestFooRepo_CreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Simulate Exec returning an error
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"INSERT INTO foos (name) VALUES ($1);",
			gomock.Any(),
		).
		Return(pgconn.NewCommandTag("INSERT 0"), errors.New("exec failed"))

	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Call under test
	rowsAffected, err := fooRepo.CreateFoo()
	require.Equal(t, int64(0), rowsAffected, "should return zero rows on error")
	require.Error(t, err)
	require.Contains(t, err.Error(), "T5O31W", "error should be wrapped with T5O31W code")
}

func TestFooRepo_DeleteFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Set expectation for the mock pgx pool.
	// Make sure the right delete is called.
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"DELETE FROM foos;",
		).
		Return(pgconn.NewCommandTag("DELETE 5"), nil)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	// Create a new foo repo
	repo := NewFooRepository(mockPool, logger)

	// Call the DeleteFoos function under test.
	rowsAffected, err := repo.DeleteFoos()
	require.NoError(t, err, "DeleteFoos should not error.")
	require.Equal(t, int64(5), rowsAffected, "should report 5 rows deleted.")
}

func TestFooRepo_DeleteFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Simulate Exec returning an error
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"DELETE FROM foos;",
		).
		Return(pgconn.NewCommandTag("DELETE 0"), errors.New("exec failed"))

	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Call under test
	rowsAffected, err := fooRepo.DeleteFoos()
	require.Equal(t, int64(0), rowsAffected, "should return zero rows on error")
	require.Error(t, err)
	require.Contains(t, err.Error(), "1BLNNL", "error should be wrapped with 1BLNNL code")
}

func TestFooRepo_UpdateFoo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Set expectation for the mock pgx pool.
	// Make sure the right update is called.
	const updateSQL = "UPDATE foos SET name = $1 WHERE id = $2;"
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			updateSQL,
			gomock.Any(),
			gomock.Any(),
		).
		Return(pgconn.NewCommandTag("UPDATE 1"), nil)

	// Create the mock logger
	logger := zaptest.NewLogger(t)

	// Create a new foo repo
	repo := NewFooRepository(mockPool, logger)

	// Call the DeleteFoos function under test.
	rowsAffected, err := repo.UpdateFoo(1)
	require.NoError(t, err, "UpdateFoo should not error.")
	require.Equal(t, int64(1), rowsAffected, "should report 1 row updated.")
}

func TestFooRepo_UpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPool(ctrl)

	// Simulate Exec returning an error for the update
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"UPDATE foos SET name = $1 WHERE id = $2;",
			gomock.Any(),
			gomock.Any(),
		).
		Return(pgconn.NewCommandTag("UPDATE 0"), errors.New("update failed"))

	logger := zaptest.NewLogger(t)
	repo := NewFooRepository(mockPool, logger)

	rowsAffected, err := repo.UpdateFoo(1)
	require.Equal(t, int64(0), rowsAffected, "should return zero rows on error")
	require.Error(t, err)
	require.Contains(t, err.Error(), "71PVZL", "error should be wrapped with 71PVZL code")
}
