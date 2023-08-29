package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	zerolog.Logger
}

func NewDefaultLogger(level string, appName, envName string) (Logger, error) {
	var logWriter io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    true,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			if ll, ok := i.(string); ok {
				return strings.ToUpper(ll)
			}
			return "????"
		},
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level (%s): %w", level, err)
	}

	// make the writer thread-safe
	logWriter = newSyncWriter(logWriter, appName, envName)

	return &defaultLogger{
		Logger: zerolog.New(logWriter).Level(logLevel).With().Timestamp().Logger(),
	}, nil
}

// MustNewDefaultLogger delegates a call NewDefaultLogger where it panics on
// error.
func MustNewDefaultLogger(level string, brokers, topic, appName, envName string) Logger {
	logger, err := NewDefaultLogger(level, appName, envName)
	if err != nil {
		panic(err)
	}

	return logger
}

func (l defaultLogger) Info(msg string, keyVals ...interface{}) {
	l.Logger.Info().Fields(getLogFields(keyVals...)).Msg(msg)
}

func (l defaultLogger) Error(msg string, keyVals ...interface{}) {
	l.Logger.Error().Fields(getLogFields(keyVals...)).Msg(msg)
}

func (l defaultLogger) Debug(msg string, keyVals ...interface{}) {
	l.Logger.Debug().Fields(getLogFields(keyVals...)).Msg(msg)
}

func (l defaultLogger) With(keyVals ...interface{}) Logger {
	return &defaultLogger{
		Logger: l.Logger.With().Fields(getLogFields(keyVals...)).Logger(),
	}
}

// OverrideWithNewLogger replaces an existing logger's internal with
// a new logger, and makes it possible to reconfigure an existing
// logger that has already been propagated to callers.
func OverrideWithNewLogger(logger Logger, level, brokers, topic, appName, envName string) error {
	ol, ok := logger.(*defaultLogger)
	if !ok {
		return fmt.Errorf("logger %T cannot be overridden", logger)
	}

	newLogger, err := NewDefaultLogger(level, appName, envName)
	if err != nil {
		return err
	}
	nl, ok := newLogger.(*defaultLogger)
	if !ok {
		return fmt.Errorf("logger %T cannot be overridden by %T", logger, newLogger)
	}

	ol.Logger = nl.Logger
	return nil
}

func getLogFields(keyVals ...interface{}) map[string]interface{} {
	if len(keyVals)%2 != 0 {
		return nil
	}

	fields := make(map[string]interface{}, len(keyVals))
	for i := 0; i < len(keyVals); i += 2 {
		fields[fmt.Sprint(keyVals[i])] = keyVals[i+1]
	}

	return fields
}
