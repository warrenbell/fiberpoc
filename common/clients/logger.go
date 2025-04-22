package clients

import (
	"os"
	"sync"

	"gitlab.com/sandstone2/fiberpoc/common/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// use github.com/pkg/errors to wrap errors with context and then return them. See https://pkg.go.dev/github.com/pkg/errors
// return errors.Wrap(err, "read failed")
// See https://github.com/uber-go/zap for Zap

// GetLogger returns a singleton instance of a configured zap.Logger.
func GetLogger() *zap.Logger {
	once.Do(func() {

		var coreFile zapcore.Core

		if *models.GlobalConfig.GetLogToFile() {
			// Set up lumberjack for log rotation.
			lumberjackLogger := &lumberjack.Logger{
				Filename:   "logs/app.json", // Log file path
				MaxSize:    5,               // Max size in megabytes before rotation
				MaxBackups: 3,               // Max number of rotated log files to retain
				MaxAge:     28,              // Max age in days to retain old logs
				Compress:   false,           // Compress rotated files or not
			}

			// File Encoder configurations.
			fileEncoderConfig := zap.NewProductionEncoderConfig()

			// Create file encoders.
			fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

			// Write syncers for file.
			fileWS := zapcore.AddSync(lumberjackLogger)

			// Create cores for file.
			coreFile = zapcore.NewCore(fileEncoder, fileWS, *models.GlobalConfig.GetLogLevel())

		}

		// Console encoder configurations.
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()

		// Create console encoders.
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		// Write syncers for console.
		consoleWS := zapcore.Lock(os.Stdout)

		// Create cores for console.
		coreConsole := zapcore.NewCore(consoleEncoder, consoleWS, *models.GlobalConfig.GetLogLevel())

		var teeCore zapcore.Core

		if *models.GlobalConfig.GetLogToFile() {
			// Combine the two cores.
			teeCore = zapcore.NewTee(coreConsole, coreFile)
		} else {
			// Use the console core only
			teeCore = zapcore.NewTee(coreConsole)
		}

		// Build the logger with caller and stack trace options.
		instance = zap.New(teeCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	})
	return instance
}
