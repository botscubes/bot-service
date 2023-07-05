package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

const (
	// Max bot title length
	MaxTitleLen = 50

	MainComponentId    = 1
	MaxPositionX       = 10000
	MaxPositionY       = 10000
	StartComponentPosX = 50
	StartComponentPosY = 50

	RedisExpire     = 1 * time.Hour
	ShutdownTimeout = 1 * time.Minute
)

type ServiceConfig struct {
	Pg            PostgresConfig
	RedisAuth     RedisAuthConfig
	Redis         RedisConfig
	WebhookBase   string `env:"WEBHOOK_BASE,required"`
	JWTKey        string `env:"JWT_SECRET_KEY,required"`
	ListenAddress string `env:"LISTEN_ADDRESS,required"`
	LoggerType    string `env:"LOGGER_TYPE,required"`
}

type PostgresConfig struct {
	Db   string `env:"POSTGRES_DB,required"`
	User string `env:"POSTGRES_USER,required"`
	Pass string `env:"POSTGRES_PASSWORD,required"`
	Host string `env:"POSTGRES_HOST,required"`
	Port string `env:"POSTGRES_PORT,required"`
}

type RedisAuthConfig struct {
	Db   int    `env:"REDIS_AUTH_DB,required"`
	Pass string `env:"REDIS_AUTH_PASS,required"`
	Host string `env:"REDIS_AUTH_HOST,required"`
	Port string `env:"REDIS_AUTH_PORT,required"`
}

type RedisConfig struct {
	Db   int    `env:"REDIS_DB,required"`
	Pass string `env:"REDIS_PASS,required"`
	Host string `env:"REDIS_HOST,required"`
	Port string `env:"REDIS_PORT,required"`
}

func GetConfig() (*ServiceConfig, error) {
	var c ServiceConfig
	if err := envconfig.Process(context.Background(), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
