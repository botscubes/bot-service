package app

// TODO: handle the panic
// WARN: bot not receive updates on app stop by panic

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"

	fastRouter "github.com/fasthttp/router"

	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	bcRedis "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/botscubes/user-service/pkg/token_storage"

	"github.com/mymmrac/telego"
)

type App struct {
	Router         *fastRouter.Router
	Server         *telego.MultiBotWebhookServer
	Bots           map[int64]*bot.TBot
	Conf           *config.ServiceConfig
	Db             *pgsql.Db
	SessionStorage token_storage.TokenStorage
	RedisAuth      *redis.Client
}

func (app *App) Run() error {
	log.Debug("App Run")

	var err error
	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.Conf, err = config.GetConfig()
	if err != nil {
		return err
	}

	app.Router = fastRouter.New()
	app.Server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{
				Handler: app.Router.Handler,
			},
			Router: app.Router,
		},
	}

	app.Bots = make(map[int64]*bot.TBot)

	app.RedisAuth = bcRedis.NewClient(&app.Conf.RedisAuth)
	app.SessionStorage = token_storage.NewRedisTokenStorage(app.RedisAuth)

	pgsqlUrl := "postgres://" + app.Conf.Pg.User + ":" + app.Conf.Pg.Pass + "@" + app.Conf.Pg.Host + ":" + app.Conf.Pg.Port + "/" + app.Conf.Pg.Db
	if app.Db, err = pgsql.OpenConnection(pgsqlUrl); err != nil {
		return err
	}

	defer app.Db.CloseConnection()

	app.addHandlers()

	go func() {
		if err = app.Server.Start(app.Conf.Bot.ListenAddress); err != nil {
			log.Error(err)
			sigs <- syscall.SIGTERM
		}
	}()

	// On close, error program
	go func() {
		<-sigs
		log.Info("Stopping...")
		for _, v := range app.Bots {
			if err := v.StopBot(true); err != nil {
				log.Error("Stop App: bot stop:\n", err)
			}
		}
		done <- struct{}{}
	}()

	log.Info("App Started")

	<-done

	return nil
}
