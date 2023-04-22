package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type ServiceConfig struct {
	Bot BotConfig
	Pg  PostgresConfig
}

type BotConfig struct {
	WebhookBase   string `env:"TBOT_WEBHOOK_BASE,required"`
	ListenAddress string `env:"TBOT_LISTEN_ADDRESS,required"`
}

type PostgresConfig struct {
	Db   string `env:"POSTGRES_DB,required"`
	User string `env:"POSTGRES_USER,required"`
	Pass string `env:"POSTGRES_PASSWORD,required"`
	Host string `env:"POSTGRES_HOST,required"`
	Port string `env:"POSTGRES_PORT,required"`
}

func GetConfig() (*ServiceConfig, error) {
	var c ServiceConfig
	if err := envconfig.Process(context.Background(), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
