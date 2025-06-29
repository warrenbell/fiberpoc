package server

import (
	"log"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/sandstone2/fiberpoc/common/clients"
	"gitlab.com/sandstone2/fiberpoc/common/models"
)

func InitServer(envFileName string) (db *clients.PgxPoolImpl, logger *zap.Logger, err error) {

	// Load environment variables from .env file if present.
	if err := godotenv.Load(envFileName); err != nil {
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

	// Get the db pool.
	db, err = clients.NewPgxPoolImpl()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error: FGI573 - Getting database connection pool.")
	}

	return db, logger, nil
}
