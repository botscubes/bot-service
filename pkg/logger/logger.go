package logger

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const envPrefix = "TBOT_"

const (
	loggerTypeVarName = "LOGGER_TYPE"
	defLogggerType    = "dev"
)

func NewLogger() (*zap.SugaredLogger, error) {
	ltype, ok := lookupEnv(loggerTypeVarName)
	if !ok {
		fmt.Printf("env %q not found, used default loggerType: %s\n", envPrefix+loggerTypeVarName, defLogggerType)
		ltype = defLogggerType
	}

	var logger *zap.Logger
	var loggerConf zap.Config
	var err error
	switch ltype {
	case "dev":
		loggerConf = zap.NewDevelopmentConfig()
	case "prod":
		loggerConf = zap.NewProductionConfig()
	default:
		return nil, errors.New("unknown logger type")
	}

	loggerConf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err = loggerConf.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(envPrefix + key)
}
