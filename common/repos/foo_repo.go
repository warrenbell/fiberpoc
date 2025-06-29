package repos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/common/interfaces"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"go.uber.org/zap"
)

/*
Install mockgen with these commands:
go get go.uber.org/mock/gomock
go install go.uber.org/mock/mockgen@latest

Then create mocks for the below interface with these commands:
mockgen \
  -destination=./mocks/mock_foo_repo.go \
  -package=mocks \
  -mock_names=FooRepoInterface=MockFooRepo \
  gitlab.com/sandstone2/fiberpoc/common/repos \
  FooRepoInterface
*/

type FooRepoInterface interface {
	GetFoos() (foos *[]models.Foo, err error)
	CreateFoo(name string) (foo *models.Foo, err error)
	DeleteFoos() (rowsAffected int64, err error)
	UpdateFoo(fooId int64, name string) (foo *models.Foo, err error)
}

type FooRepo struct {
	db     *interfaces.PgxPoolInterface
	logger *zap.Logger
}

func NewFooRepository(db interfaces.PgxPoolInterface, logger *zap.Logger) *FooRepo {
	return &FooRepo{db: &db, logger: logger}
}

func (fooRepo *FooRepo) GetFoos() (foos *[]models.Foo, err error) {
	foos = &[]models.Foo{}

	rows, err := (*fooRepo.db).Query(context.Background(), "SELECT id, name FROM foos ORDER BY id;")
	if err != nil {
		return nil, errors.Wrap(err, "Error: 30UUBR - Quering foos from db. Error")
	}
	defer rows.Close()

	for rows.Next() {
		foo := models.Foo{}

		if err := rows.Scan(&foo.ID, &foo.Name); err != nil {
			return nil, errors.Wrap(err, "Error: YN80XB - Scanning row of foos from db.")
		}

		*foos = append(*foos, foo)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "Error: XV4HHL - Processing rows of foos from db.")
	}

	return foos, nil
}

func (fooRepo *FooRepo) CreateFoo(name string) (foo *models.Foo, err error) {
	foo = &models.Foo{}
	err = (*fooRepo.db).QueryRow(
		context.Background(),
		"INSERT INTO foos (name) VALUES ($1) RETURNING id, name;",
		name,
	).Scan(&foo.ID, &foo.Name)

	if err != nil {
		return nil, errors.Wrap(err, "Error: WOPUDO - Inserting foo into database.")
	}

	return foo, nil
}
func (fooRepo *FooRepo) DeleteFoos() (rowsAffected int64, err error) {
	var result pgconn.CommandTag

	result, err = (*fooRepo.db).Exec(context.Background(), "DELETE FROM foos;")
	if err != nil {
		return 0, errors.Wrap(err, "Error: 1BLNNL - Deleteing foos from database.")
	}

	return result.RowsAffected(), nil
}

func (fooRepo *FooRepo) UpdateFoo(fooId int64, name string) (foo *models.Foo, err error) {
	foo = &models.Foo{}
	err = (*fooRepo.db).QueryRow(
		context.Background(),
		"UPDATE foos SET name = $1 WHERE id = $2 RETURNING id, name;",
		name,
		fooId,
	).Scan(&foo.ID, &foo.Name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, fmt.Sprintf("Error: BATWXG - No foo found with given ID: %d", fooId))
		}
		return nil, errors.Wrap(err, "Error: 2H6YX9 - Updating foo in database.")
	}

	return foo, nil
}
