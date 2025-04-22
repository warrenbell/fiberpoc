package interfaces

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
)

// PgxRowInterface is an interface for a single-row result (e.g., QueryRow).
type PgxRowInterface interface {
	Scan(dest ...interface{}) error
}

// PgxRowsInterface is an interface for multi-row results (e.g., Query).
type PgxRowsInterface interface {
	Close()
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

// PgxPoolInterface is our main interface that wraps the methods we need from pgxpool.Pool.
type PgxPoolInterface interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (PgxRowsInterface, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) PgxRowInterface
	Close()
}

/*
Install mockgen with these commands:
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

Then create mocks for the above interfaces with these commands:
mockgen \
  -destination=./mocks/mock_pgx_row.go \
  -package=mocks \
  -mock_names=PgxRowInterface=MockPgxRow \
  gitlab.com/sandstone2/fiberpoc/common/interfaces \
  PgxRowInterface

 mockgen \
  -destination=./mocks/mock_pgx_rows.go \
  -package=mocks \
  -mock_names=PgxRowsInterface=MockPgxRows \
  gitlab.com/sandstone2/fiberpoc/common/interfaces \
  PgxRowsInterface

mockgen \
  -destination=./mocks/mock_pgx_pool.go \
  -package=mocks \
  -mock_names=PgxPoolInterface=MockPgxPool \
  gitlab.com/sandstone2/fiberpoc/common/interfaces \
  PgxPoolInterface
*/
