package repos

import (
	"context"

	"github.com/Pallinder/go-randomdata"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/common/interfaces"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"go.uber.org/zap"
)

/*
Install mockgen with these commands:
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

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
	CreateFoo() (rowsAffected int64, err error)
	DeleteFoos() (rowsAffected int64, err error)
	UpdateFoo(fooId int64) (rowsAffected int64, err error)
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

	rows, err := (*fooRepo.db).Query(context.Background(), "SELECT * FROM foos;")
	if err != nil {
		return nil, errors.Wrap(err, "Error: 30UUBR - Quering foos from db.")
	}
	defer rows.Close()

	for rows.Next() {
		var foo models.Foo

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

func (fooRepo *FooRepo) CreateFoo() (rowsAffected int64, err error) {
	name := randomdata.SillyName()

	var result pgconn.CommandTag

	result, err = (*fooRepo.db).Exec(context.Background(), "INSERT INTO foos (name) VALUES ($1);", name)
	if err != nil {
		return 0, errors.Wrap(err, "Error: T5O31W - Inserting foo into database.")
	}

	return result.RowsAffected(), nil
}

func (fooRepo *FooRepo) DeleteFoos() (rowsAffected int64, err error) {
	var result pgconn.CommandTag

	result, err = (*fooRepo.db).Exec(context.Background(), "DELETE FROM foos;")
	if err != nil {
		return 0, errors.Wrap(err, "Error: 1BLNNL - Deleteing foos from database.")
	}

	return result.RowsAffected(), nil
}

func (fooRepo *FooRepo) UpdateFoo(fooId int64) (rowsAffected int64, err error) {
	name := randomdata.SillyName()

	var result pgconn.CommandTag

	result, err = (*fooRepo.db).Exec(context.Background(), "UPDATE foos SET name = $1 WHERE id = $2;", name, fooId)
	if err != nil {
		return 0, errors.Wrap(err, "Error: 71PVZL - Updating foo in database.")
	}

	return result.RowsAffected(), nil
}
