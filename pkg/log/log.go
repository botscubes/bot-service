package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const envPrefix = "TBOT_"

const (
	logLevelVarName  = "LOG_LEVEL"
	logFormatVarName = "LOG_FORMAT"
	logPathVarName   = "LOG_PATH"
	defaultLevel     = logrus.InfoLevel
)

var log = newLogger()

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(envPrefix + key)
}

func newLogger() *logrus.Logger {
	level := getLevel()
	formatter := getFormatter()
	output := getOutput()

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	logger.SetOutput(output)
	return logger
}

func getLevel() logrus.Level {
	lvl, ok := lookupEnv(logLevelVarName)
	if !ok {
		fmt.Printf("env %q not found, used default log level\n", logLevelVarName)
		return defaultLevel
	}

	var level logrus.Level
	switch strings.ToLower(lvl) {
	case "trace":
		level = logrus.TraceLevel
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	default:
		fmt.Println("unknown log level ", lvl)
		level = defaultLevel
	}

	return level
}

func getFormatter() logrus.Formatter {
	formatName, _ := lookupEnv(logFormatVarName)
	switch strings.ToLower(formatName) {
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006.01.02 15:04:05",
			DisableLevelTruncation: true,
			PadLevelText:           true,
		}
	}
}

// setOutput sets the logger output.
// Trying to open a log file.
// For output to stderr, do not set the value of the environment variable of the log file path (see logPathVarName)
func getOutput() io.Writer {
	filePath, ok := lookupEnv(logPathVarName)
	if !ok {
		return os.Stderr
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("failed open log file %q;\n%s", filePath, err)
	}
	return file
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Print(args ...interface{}) {
	log.Print(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warning(args ...interface{}) {
	log.Warning(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}
