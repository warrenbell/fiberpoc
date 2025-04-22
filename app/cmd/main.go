package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/sandstone2/fiberpoc/app/handlers"
	"gitlab.com/sandstone2/fiberpoc/common/clients"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"gitlab.com/sandstone2/fiberpoc/common/repos"
	"gitlab.com/sandstone2/fiberpoc/common/services"
)

func main() {
	// Initialize the server.
	db, logger, err := initServer()
	if err != nil {
		// Can not use logger here.
		log.Fatalf("Error: LBTF9J - Initializing the server. Error: %v", err)
	}

	// Flush out the logger on server exit.
	defer logger.Sync()

	// Close the db pool on server exit.
	defer db.Close()

	// Inject all dependencies.
	fooRepo := repos.NewFooRepository(db, logger)
	fooService := services.NewFooService(fooRepo, logger)
	fooHandler := handlers.NewFooHandler(fooService, logger)

	// Create the Fiber app.
	app := fiber.New()

	// Create the routes.
	app.Get("/foos", fooHandler.HandleGetFoos)
	app.Post("/foo", fooHandler.HandleCreateFoo)
	app.Delete("/foos", fooHandler.HandleDeleteFoos)
	app.Patch("/foo", fooHandler.HandleUpdateFoo)

	// Start the Fiber server in a separate goroutine.
	go func(app *fiber.App) {
		clients.GetLogger().Info("Fiber listening on port :3000")
		if err := app.Listen(":3000"); err != nil {
			clients.GetLogger().Sugar().Fatalf("Error: L4AXAX - Starting Fiber server. Error: %v", err)
		}
	}(app)

	// Handle graceful shutdown signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Block until we get a signal.

	clients.GetLogger().Info("Shutting down Fiber server...")

	// Shutdown in a go routine so we can time out if needed.
	done := make(chan error, 1)
	go func(app *fiber.App) {
		done <- app.Shutdown()
	}(app)

	// Wait for either shutdown to complete or time out in 5 seconds
	select {
	case err := <-done:
		if err != nil {
			clients.GetLogger().Sugar().Fatalf("Error: B50QRT - Shutting down Fiber server. Error: %v", err)
		} else {
			clients.GetLogger().Info("Fiber server shutdown complete.")
		}
	case <-time.After(5 * time.Second):
		log.Println("Graceful shutdown for Fiber server timed out. Forcing shutdown.")
	}
}

func initServer() (db *clients.PgxPoolImpl, logger *zap.Logger, err error) {

	// Load environment variables from .env file if present.
	if err := godotenv.Load(); err != nil {
		// we use log here instead of Zap because the env vars are not loaded yet
		log.Println("No .env file found or unable to load, proceeding with environment variables already in the environment")
	}

	// Parse and validate environment variables into a config struct.
	config := models.AppConfig{}
	if err := env.Parse(&config); err != nil {
		// Wrap the error
		return nil, nil, errors.Wrap(err, "Error: YN80XB - Parsing and validating env vars")
	}

	models.GlobalConfig = &config

	logger = clients.GetLogger()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error: O6FSH5 - Getting logger.")
	}

	// Get the db pool.
	db, err = clients.NewPgxPoolImpl()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error: FGI573 - Getting database connection pool.")
	}

	return db, logger, nil
}
