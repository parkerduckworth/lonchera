// Package log is a wrapper package for a [logrus](https://github.com/sirupsen/logrus)
// instance, providing structured, leveled logging throughout the application. The log
// level is set in the application configuration YAML files.
package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/parkerduckworth/lonchera/app/config"
	"github.com/sirupsen/logrus"
)

var logger Logger

type Logger struct {
	internalLogger *logrus.Entry
}

func Trace(args ...interface{}) {
	logger.internalLogger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.internalLogger.Tracef(format, args...)
}

func Debug(args ...interface{}) {
	logger.internalLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.internalLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.internalLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.internalLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	logger.internalLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.internalLogger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logger.internalLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.internalLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.internalLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.internalLogger.Fatalf(format, args...)
}

func (logger *Logger) SetLevel(level string) {
	switch level {
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

}

// Setup defines the logger configuration and log message structure
// and sets the level of the logger
func Setup() {
	logrus.SetReportCaller(false)
	formatter := &logrus.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmed everything before the file name
			// and added a line number in the end
			return "-", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
	lgr := logrus.WithContext(context.Background())

	logger = Logger{
		internalLogger: lgr,
	}
	logger.SetLevel(config.Conf.Logger.Level)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
