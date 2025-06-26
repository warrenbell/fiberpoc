package models

import "go.uber.org/zap/zapcore"

var GlobalConfig Config

type Config interface {
	GetLogLevel() *zapcore.Level
	GetLogToFile() *bool
	GetPostgresUrl() *string
	GetGoogleOidcClientId() *string
	GetGoogleOidcClientSecret() *string
	GetGoogleOidcProviderUrl() *string
	GetRedirectUri() *string
}

type AppConfig struct {
	PostgresUrl            string        `env:"POSTGRESQL_URL,required"`
	LogLevel               zapcore.Level `env:"LOG_LEVEL" envDefault:"debug"`
	LogToFile              bool          `env:"LOG_TO_FILE" envDefault:"false"`
	GoogleOidcClientId     string        `env:"GOOGLE_OIDC_CLIENT_ID,required"`
	GoogleOidcClientSecret string        `env:"GOOGLE_OIDC_CLIENT_SECRET,required"`
	GoogleOidcProviderUrl  string        `env:"GOOGLE_OIDC_PROVIDER_URL,required"`
	RedirectUri            string        `env:"REDIRECT_URI,required"`
}

func (appConfig *AppConfig) GetPostgresUrl() *string {
	return &appConfig.PostgresUrl
}

func (appConfig *AppConfig) GetLogLevel() *zapcore.Level {
	return &appConfig.LogLevel
}

func (appConfig *AppConfig) GetLogToFile() *bool {
	return &appConfig.LogToFile
}

func (appConfig *AppConfig) GetGoogleOidcClientId() *string {
	return &appConfig.GoogleOidcClientId
}

func (appConfig *AppConfig) GetGoogleOidcClientSecret() *string {
	return &appConfig.GoogleOidcClientSecret
}

func (appConfig *AppConfig) GetGoogleOidcProviderUrl() *string {
	return &appConfig.GoogleOidcProviderUrl
}

func (appConfig *AppConfig) GetRedirectUri() *string {
	return &appConfig.RedirectUri
}
