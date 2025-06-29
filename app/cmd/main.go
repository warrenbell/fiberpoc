package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"gitlab.com/sandstone2/fiberpoc/app/handlers"
	"gitlab.com/sandstone2/fiberpoc/app/middleware"
	"gitlab.com/sandstone2/fiberpoc/app/server"

	"gitlab.com/sandstone2/fiberpoc/common/clients"
	"gitlab.com/sandstone2/fiberpoc/common/repos"
	"gitlab.com/sandstone2/fiberpoc/common/services"
)

func main() {
	// Initialize the server.
	db, logger, err := server.InitServer(".env")
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

	authcService, err := services.NewAuthcService(logger)
	if err != nil {
		logger.Sugar().Fatalf("Error: A18S5B - Creating AuthcService. Error: %v", err)
	}
	authcHandler := handlers.NewAuthcHandler(authcService, logger)

	engine := html.New("./templates", ".html")
	engine.Reload(true)
	// Create the Fiber app.
	app := fiber.New(fiber.Config{Views: engine})

	// Create the routes.

	app.Get("/", authcHandler.HandleRoot)
	app.Get("/login", authcHandler.HandleLogin)
	app.Get("/callback", authcHandler.HandleOauthCallback)
	app.Get("/foos", middleware.AuthcMiddleware(authcService.GetVerifier(), logger), fooHandler.HandleGetFoos)
	app.Post("/foos", middleware.AuthcMiddleware(authcService.GetVerifier(), logger), fooHandler.HandleCreateFoo)
	app.Delete("/foos", middleware.AuthcMiddleware(authcService.GetVerifier(), logger), fooHandler.HandleDeleteFoos)
	app.Put("/foos/:id", middleware.AuthcMiddleware(authcService.GetVerifier(), logger), fooHandler.HandleUpdateFoo) // Replace all fields with new ones.

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
