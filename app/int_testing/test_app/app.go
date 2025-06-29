package testapp

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/app/handlers"
	"gitlab.com/sandstone2/fiberpoc/app/server"
	"gitlab.com/sandstone2/fiberpoc/common/clients"
	"gitlab.com/sandstone2/fiberpoc/common/repos"
	"gitlab.com/sandstone2/fiberpoc/common/services"
	"go.uber.org/zap"
)

var db *clients.PgxPoolImpl
var logger *zap.Logger

func GetApp() (app *fiber.App, err error) {
	db, logger, err = server.InitServer(".env.tst")
	if err != nil {
		return nil, errors.Wrap(err, "Error: LBTF9J - Initializing the server. Error: %v")
	}

	// Inject all dependencies.
	fooRepo := repos.NewFooRepository(db, logger)
	fooService := services.NewFooService(fooRepo, logger)
	fooHandler := handlers.NewFooHandler(fooService, logger)

	// Create the Fiber app.
	app = fiber.New(fiber.Config{})

	app.Get("/foos", fooHandler.HandleGetFoos)

	return app, nil
}

func CloseDbAndLogger() {
	// Flush out the logger on server exit.
	logger.Sync()

	// Close the db pool on server exit.
	db.Close()
}
