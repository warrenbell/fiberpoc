package clients

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/common/interfaces"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

// pgxPoolImpl is our production wrapper around *pgxpool.Pool.
type PgxPoolImpl struct {
	pool *pgxpool.Pool
}

// NewPgxPoolImpl constructs a Pool from a real *pgxpool.Pool.
func NewPgxPoolImpl() (*PgxPoolImpl, error) {
	url := *models.GlobalConfig.GetPostgresUrl()

	connConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		// Fatal kill server
		return nil, errors.Wrap(err, "Error: V78BO4 - Parsing the database configs from the url.")
	}

	var pool *pgxpool.Pool

	pool, err = pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		// Fatal kill server
		return nil, errors.Wrap(err, "Error: GKK2EB - Creating the database pool.")
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		// Fatal kill server
		return nil, errors.Wrap(err, "Error: HT9IXG - Testing the database pool.")
	}

	GetLogger().Sugar().Info("Database is Connected")

	return &PgxPoolImpl{pool: pool}, nil
}

// Exec delegates to the real pool.Exec.
func (p *PgxPoolImpl) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

// Query delegates to the real pool.Query, returning a *pgxRows wrapper.
func (p *PgxPoolImpl) Query(ctx context.Context, sql string, args ...interface{}) (interfaces.PgxRowsInterface, error) {
	rawRows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxRows{rows: rawRows}, nil
}

// QueryRow delegates to the real pool.QueryRow, returning a *pgxRow wrapper.
func (p *PgxPoolImpl) QueryRow(ctx context.Context, sql string, args ...interface{}) interfaces.PgxRowInterface {
	return &pgxRow{row: p.pool.QueryRow(ctx, sql, args...)}
}

func (p *PgxPoolImpl) Close() {
	p.pool.Close()
}

// pgxRows wraps pgx.Rows to implement our Rows interface.
type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Close() {
	r.rows.Close()
}

func (r *pgxRows) Next() bool {
	return r.rows.Next()
}

func (r *pgxRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *pgxRows) Err() error {
	return r.rows.Err()
}

// pgxRow wraps pgx.Row to implement our Row interface.
type pgxRow struct {
	row pgx.Row
}

func (r *pgxRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}
