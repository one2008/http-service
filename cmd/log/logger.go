package log

import (
	"io"
	"sync"
)

const (
	// log级别
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type Logger interface {
	Debug(msg string, keyVals ...interface{})
	Info(msg string, keyVals ...interface{})
	Error(msg string, keyVals ...interface{})

	With(keyVals ...interface{}) Logger
}

type syncWriter struct {
	sync.Mutex
	io.Writer
	topic   string
	appName string
	envName string
}

func newSyncWriter(w io.Writer, appName, envName string) io.Writer {

	return &syncWriter{Writer: w, appName: appName, envName: envName}
}

func (w *syncWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()

	return w.Writer.Write(p)
}

type KafkaLog struct {
	AppName string `json:"app_name"`
	App     string `json:"app"`
	Env     string `json:"env"`
	EnvName string `json:"env_name"`
	Message string `json:"message"`
}
