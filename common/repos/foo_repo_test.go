// myapp/repository/repository_test.go
package repos

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
	const expectedQuery = "SELECT id, name FROM foos ORDER BY id;"
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
		Query(gomock.Any(), "SELECT id, name FROM foos ORDER BY id;").
		Return(mockRows, errors.New("query failed"))

	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Call under test
	_, err := fooRepo.GetFoos()
	require.Error(t, err)
	require.Contains(t, err.Error(), "30UUBR", "error should be wrapped with 30UUBR code")

	// 2) Test mockRows.Scan failed
	mockPool.EXPECT().
		Query(gomock.Any(), "SELECT id, name FROM foos ORDER BY id;").
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
		Query(gomock.Any(), "SELECT id, name FROM foos ORDER BY id;").
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

	// Create the mock pgx pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Create the mock Row
	mockRow := mocks.NewMockPgxRow(ctrl)

	// Set up the expected query
	mockPool.EXPECT().
		QueryRow(
			gomock.Any(),
			"INSERT INTO foos (name) VALUES ($1) RETURNING id, name;",
			"Test Foo",
		).
		Return(mockRow)

	// Mock the scan to populate our expected result
	mockRow.EXPECT().
		Scan(gomock.Any()).
		DoAndReturn(func(dest ...any) error {
			id := dest[0].(*int)
			name := dest[1].(*string)
			*id = 1
			*name = "Test Foo"
			return nil
		})

	// Create logger and FooRepo
	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Act
	foo, err := fooRepo.CreateFoo("Test Foo")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, foo)
	require.Equal(t, 1, foo.ID)
	require.Equal(t, "Test Foo", foo.Name)
}

func TestFooRepo_CreateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pgx pool
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Simulate QueryRow().Scan() returning an error
	mockRow := mocks.NewMockPgxRow(ctrl)
	mockPool.EXPECT().
		QueryRow(
			gomock.Any(),
			"INSERT INTO foos (name) VALUES ($1) RETURNING id, name;",
			"Bad Foo",
		).
		Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		Return(errors.New("scan failed"))

	logger := zaptest.NewLogger(t)
	fooRepo := NewFooRepository(mockPool, logger)

	// Call under test
	foo, err := fooRepo.CreateFoo("Bad Foo")

	require.Nil(t, foo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "WOPUDO", "should wrap with correct error code")
}

func TestFooRepo_DeleteFoos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool interface
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Expect the DELETE SQL call
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"DELETE FROM foos;",
		).
		Return(pgconn.NewCommandTag("DELETE 5"), nil)

	logger := zaptest.NewLogger(t)
	repo := NewFooRepository(mockPool, logger)

	// Act
	rowsAffected, err := repo.DeleteFoos()

	// Assert
	require.NoError(t, err, "DeleteFoos should not return error")
	require.Equal(t, int64(5), rowsAffected, "should report 5 rows deleted")
}

func TestFooRepo_DeleteFoos_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock pool interface
	mockPool := mocks.NewMockPgxPool(ctrl)

	// Simulate Exec returning an error
	mockPool.
		EXPECT().
		Exec(
			gomock.Any(),
			"DELETE FROM foos;",
		).
		Return(pgconn.CommandTag{}, errors.New("exec failed"))

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

	// Create mocks
	mockPool := mocks.NewMockPgxPool(ctrl)
	mockRow := mocks.NewMockPgxRow(ctrl)

	// Set expected query
	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			"UPDATE foos SET name = $1 WHERE id = $2 RETURNING id, name;",
			"Updated Foo",
			int64(1),
		).
		Return(mockRow)

	// Simulate Scan populating values
	mockRow.EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		DoAndReturn(func(dest ...any) error {
			*(dest[0].(*int)) = 1
			*(dest[1].(*string)) = "Updated Foo"
			return nil
		})

	logger := zaptest.NewLogger(t)
	repo := NewFooRepository(mockPool, logger)

	// Act
	foo, err := repo.UpdateFoo(1, "Updated Foo")

	// Assert
	require.NoError(t, err)
	require.NotNil(t, foo)
	require.Equal(t, 1, foo.ID)
	require.Equal(t, "Updated Foo", foo.Name)
}

func TestFooRepo_UpdateFoo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := mocks.NewMockPgxPool(ctrl)
	mockRow := mocks.NewMockPgxRow(ctrl)

	// Simulate QueryRow().Scan() returning an error
	mockPool.
		EXPECT().
		QueryRow(
			gomock.Any(),
			"UPDATE foos SET name = $1 WHERE id = $2 RETURNING id, name;",
			"Bad Name",
			int64(99),
		).
		Return(mockRow)

	mockRow.
		EXPECT().
		Scan(gomock.Any(), gomock.Any()).
		Return(errors.New("update failed"))

	logger := zaptest.NewLogger(t)
	repo := NewFooRepository(mockPool, logger)

	// Act
	foo, err := repo.UpdateFoo(99, "Bad Name")

	// Assert
	require.Nil(t, foo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "2H6YX9", "error should be wrapped with 2H6YX9 code")
}
