package app

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
	Bots           map[string]*bot.TBot
	Conf           *config.ServiceConfig
	Db             *pgsql.Db
	SessionStorage token_storage.TokenStorage
	RedisAuth      *redis.Client
}

const envPrefix = "TBOT_"

func New() *App {
	var app App

	var err error
	app.Conf, err = config.GetConfig()
	if err != nil {
		log.Fatal("GetConfig:\n", err)
	}

	app.Router = fastRouter.New()
	app.Server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{
				Handler: nil,
			},
			Router: app.Router,
		},
	}

	app.Bots = make(map[string]*bot.TBot)

	app.RedisAuth = bcRedis.NewClient(&app.Conf.RedisAuth)
	app.SessionStorage = token_storage.NewRedisTokenStorage(app.RedisAuth)

	return &app
}

func (app *App) Run() {
	log.Debug("App Run")

	var err error

	pgsqlUrl := "postgres://" + app.Conf.Pg.User + ":" + app.Conf.Pg.Pass + "@" + app.Conf.Pg.Host + ":" + app.Conf.Pg.Port + "/" + app.Conf.Pg.Db
	if app.Db, err = pgsql.OpenConnection(pgsqlUrl); err != nil {
		log.Error("Connection Postgresql error ", err)
	}
	defer app.Db.CloseConnection()

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = app.Server.Start(app.Conf.Bot.ListenAddress); err != nil {
			log.Error("Start server: \n", err)
		}
	}()

	botToken, ok := env(envPrefix + "TOKEN")
	assert(ok, "Environment variable "+envPrefix+"TOKEN not found")

	app.Bots[botToken] = new(bot.TBot)
	app.Bots[botToken].Bot, _ = bot.NewBot(&botToken)

	if err = app.Bots[botToken].StartBot(app.Conf.Bot.WebhookBase, app.Conf.Bot.ListenAddress, app.Server); err != nil {
		log.Error("Start bot\n", err)
	}

	// On exit, close, error program
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
	log.Info("App Done")
}

func env(name string) (string, bool) {
	return os.LookupEnv(name)
}

// Check ok and exit program if ok is false
func assert(ok bool, args ...any) {
	if !ok {
		log.Fatal(append([]any{"FATAL:"}, args...)...)
	}
}
