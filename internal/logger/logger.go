package logger

import (
	log "github.com/sirupsen/logrus"
	"path"
	"runtime"
	"ypeskov/go-password-manager/internal/config"
)

type Logger struct {
	*log.Logger
}

type ContextHook struct{}

func (hook ContextHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook ContextHook) Fire(entry *log.Entry) error {
	if pc, file, line, ok := runtime.Caller(8); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["file"] = path.Base(file)
		entry.Data["func"] = path.Base(funcName)
		entry.Data["line"] = line
	}

	return nil
}

func New(cfg *config.Config) *Logger {
	l := log.New()

	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Warnf("Invalid log level '%s'. Using 'info' level as default.", cfg.LogLevel)
		level = log.InfoLevel
	}
	l.SetLevel(level)
	l.AddHook(ContextHook{})

	return &Logger{l}
}
