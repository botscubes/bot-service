package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
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
	var err error
	switch ltype {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func lookupEnv(key string) (string, bool) {
	return os.LookupEnv(envPrefix + key)
}
