package services

import (
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"gitlab.com/sandstone2/fiberpoc/common/repos"
	"go.uber.org/zap"
)

/*
Install mockgen with these commands:
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

Then create mocks for the below interface with these commands:
mockgen \
  -destination=./mocks/mock_foo_service.go \
  -package=mocks \
  -mock_names=FooServiceInterface=MockFooService \
  gitlab.com/sandstone2/fiberpoc/common/services \
  FooServiceInterface
*/

type FooServiceInterface interface {
	GetFoos() (foos *[]models.Foo, err error)
	CreateFoo() (rowsAffected int64, err error)
	DeleteFoos() (rowsAffected int64, err error)
	UpdateFoo(fooId int64) (rowsAffected int64, err error)
}

type FooService struct {
	fooRepo *repos.FooRepoInterface
	logger  *zap.Logger
}

func NewFooService(fooRepo repos.FooRepoInterface, logger *zap.Logger) *FooService {
	return &FooService{fooRepo: &fooRepo, logger: logger}
}

func (fooService *FooService) GetFoos() (foos *[]models.Foo, err error) {
	foos, err = (*fooService.fooRepo).GetFoos()
	if err != nil {
		return nil, errors.Wrap(err, "Error: WZDCXT - Getting foos.")
	}
	return foos, nil
}

func (fooService *FooService) CreateFoo() (rowsAffected int64, err error) {
	rowsAffected, err = (*fooService.fooRepo).CreateFoo()
	if err != nil {
		return 0, errors.Wrap(err, "Error: DWA4G7 - Creating foos.")
	}

	return rowsAffected, nil
}

func (fooService *FooService) DeleteFoos() (rowsAffected int64, err error) {
	rowsAffected, err = (*fooService.fooRepo).DeleteFoos()
	if err != nil {
		return 0, errors.Wrap(err, "Error: BA8TAX - Deleting foos.")
	}

	return rowsAffected, nil
}

func (fooService *FooService) UpdateFoo(fooId int64) (rowsAffected int64, err error) {
	rowsAffected, err = (*fooService.fooRepo).UpdateFoo(fooId)
	if err != nil {
		return 0, errors.Wrap(err, "Error: GZNHKW - Updating foos.")
	}

	return rowsAffected, nil
}
