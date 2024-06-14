package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

const (
	MainComponentId    = 1
	StartComponentPosX = 100
	StartComponentPosY = 100

	RedisExpire     = 1 * time.Hour
	ShutdownTimeout = 1 * time.Minute
	NatsReqTimeout  = 5 * time.Second
)

type ServiceConfig struct {
	Pg            PostgresConfig
	RedisAuth     RedisAuthConfig
	Redis         RedisConfig
	WebhookDomain string `env:"WEBHOOK_DOMAIN,required"`
	WebhookPath   string `env:"WEBHOOK_PATH,required"`
	JWTKey        string `env:"JWT_SECRET_KEY,required"`
	ListenAddress string `env:"LISTEN_ADDRESS,required"`
	LoggerType    string `env:"LOGGER_TYPE,required"`
	NatsURL       string `env:"NATS_URL,required"`
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
