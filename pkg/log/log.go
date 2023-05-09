package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
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
	output := setOutput()

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
		fmt.Printf("unknown log level %q", lvl)
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
func setOutput() io.Writer {
	filePath, ok := lookupEnv(logPathVarName)
	if !ok {
		return os.Stderr
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) //nolint:gosec,gomnd
	if err != nil {
		fmt.Printf("failed open log file %q;\n%s", filePath, err)
	}
	return file
}

func callingLine() string {
	const skip = 2
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf("\nCalled from %s, line #%d, func: %v\n", file, line, runtime.FuncForPC(pc).Name())
	}
	return ""
}

func Debugf(format string, args ...any) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	log.Infof(format, args...)
}

func Printf(format string, args ...any) {
	log.Printf(format, args...)
}

func Warnf(format string, args ...any) {
	log.Warnf(format, args...)
}

func Warningf(format string, args ...any) {
	log.Warningf(format, args...)
}

func Errorf(format string, args ...any) {
	format += callingLine()
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	format += callingLine()
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...any) {
	log.Panicf(format, args...)
}

func Debug(args ...any) {
	log.Debug(args...)
}

func Info(args ...any) {
	log.Info(args...)
}

func Print(args ...any) {
	log.Print(args...)
}

func Warn(args ...any) {
	log.Warn(args...)
}

func Warning(args ...any) {
	log.Warning(args...)
}

func Error(args ...any) {
	args = append(args, callingLine())
	log.Error(args...)
}

func Fatal(args ...any) {
	args = append(args, callingLine())
	log.Fatal(args...)
}

func Panic(args ...any) {
	log.Panic(args...)
}
