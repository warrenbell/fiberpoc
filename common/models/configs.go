package models

import "go.uber.org/zap/zapcore"

var GlobalConfig Config

type Config interface {
	GetLogLevel() *zapcore.Level
	GetLogToFile() *bool
	GetPostgresUrl() *string
}

type AppConfig struct {
	postgresUrl string        `env:"POSTGRESQL_URL,required"`
	LogLevel    zapcore.Level `env:"LOG_LEVEL" envDefault:"debug"`
	LogToFile   bool          `env:"LOG_TO_FILE" envDefault:"false"`
}

func (appConfig *AppConfig) GetPostgresUrl() *string {
	return &appConfig.postgresUrl
}

func (appConfig *AppConfig) GetLogLevel() *zapcore.Level {
	return &appConfig.LogLevel
}

func (appConfig *AppConfig) GetLogToFile() *bool {
	return &appConfig.LogToFile
}
