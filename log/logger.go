package log

import (
	"github.com/op/go-logging"
	"os"
)

const (
	DEBUG   = "DEBUG"
	INFO    = "INFO"
	WARNING = "WARNING"
	NOTICE  = "NOTICE"
	ERROR   = "ERROR"

	BoltModule   = "bolt"
)

var Bolt *logging.Logger

var (
	defaultBackend logging.Backend
	defaultFormat  = "%{color}[%{module}][%{level:.5s}] ▶ %{time:15:04:05.000} %{shortfile} %{message} %{color:reset}"
	loggerManager  *LoggerManager
)

type LoggerManager struct {
	loggers map[string]*logging.Logger
}

func NewLoggerManager() *LoggerManager {
	return &LoggerManager{
		loggers: make(map[string]*logging.Logger),
	}
}

func init() {
	loggerBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logConsoleFormat := logging.MustStringFormatter(defaultFormat)
	loggerFormatteredBackend := logging.NewBackendFormatter(loggerBackend, logConsoleFormat)
	logging.SetBackend(loggerFormatteredBackend)
	defaultBackend = loggerFormatteredBackend

	// init logger manager
	loggerManager = NewLoggerManager()
	Bolt = GetLogger(BoltModule, DEBUG)
}

func newBackendLeveled() logging.LeveledBackend {
	return logging.AddModuleLevel(defaultBackend)
}

func GetLogger(module, level string) *logging.Logger {
	if logger, exist := loggerManager.loggers[module]; exist {
		return logger
	}
	l, _ := logging.LogLevel(level)
	logger := logging.MustGetLogger(module)
	bl := newBackendLeveled()
	bl.SetLevel(l, module)
	logger.SetBackend(bl)
	loggerManager.loggers[module] = logger
	return logger
}
