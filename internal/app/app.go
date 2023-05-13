package app

// WARN: bot not receive updates on app stop by panic

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	fastRouter "github.com/fasthttp/router"

	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/database/redisauth"
	"github.com/botscubes/user-service/pkg/token_storage"

	"github.com/mymmrac/telego"
)

type App struct {
	Router         *fastRouter.Router
	Server         *telego.MultiBotWebhookServer
	BotService     *bot.BotService
	Conf           *config.ServiceConfig
	Db             *pgsql.Db
	SessionStorage token_storage.TokenStorage
	RedisAuth      *redis.Client
	Redis          *rdb.Rdb
	Log            *zap.SugaredLogger
}

func (app *App) Run(logger *zap.SugaredLogger) error {
	var err error
	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.Conf, err = config.GetConfig()
	if err != nil {
		return err
	}

	app.Log = logger

	app.Router = fastRouter.New()
	app.Server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{
				Handler: app.Router.Handler,
			},
			Router: app.Router,
		},
	}

	app.BotService = bot.NewBotService(&app.Conf.Bot, app.Server)

	app.RedisAuth = redisauth.NewClient(&app.Conf.RedisAuth)
	app.SessionStorage = token_storage.NewRedisTokenStorage(app.RedisAuth)

	app.Redis = rdb.NewClient(&app.Conf.Redis)

	pgsqlUrl := "postgres://" + app.Conf.Pg.User + ":" + app.Conf.Pg.Pass + "@" + app.Conf.Pg.Host + ":" + app.Conf.Pg.Port + "/" + app.Conf.Pg.Db
	if app.Db, err = pgsql.OpenConnection(pgsqlUrl); err != nil {
		return err
	}

	defer app.Db.CloseConnection()

	app.regiterHandlers()

	go func() {
		if err = app.Server.Start(app.Conf.Bot.ListenAddress); err != nil {
			app.Log.Error(err)
			sigs <- syscall.SIGTERM
		}
	}()

	// On close, error program
	go func() {
		<-sigs
		app.Log.Info("Stopping...")
		if err := app.BotService.StopBots(); err != nil {
			app.Log.Info("bots stop:\n", err)
		}
		done <- struct{}{}
	}()

	app.Log.Info("App Started")

	<-done

	return nil
}
