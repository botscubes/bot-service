package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"

	fastRouter "github.com/fasthttp/router"

	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
)

type App struct {
	router *fastRouter.Router
	server *telego.MultiBotWebhookServer
	bots   map[string]*bot.TBot
	conf   *config.ServiceConfig
	db     *pgsql.Db
}

var app App

const envPrefix = "TBOT_"

func Run() {
	log.Debug("Init")

	var err error = nil

	botToken, ok := env(envPrefix + "TOKEN")
	assert(ok, "Environment variable "+envPrefix+"TOKEN not found")

	app.conf, err = config.GetConfig()
	if err != nil {
		log.Fatal("GetConfig:\n", err)
	}

	pgsqlUrl := "postgres://" + app.conf.Pg.User + ":" + app.conf.Pg.Pass + "@" + app.conf.Pg.Host + ":" + app.conf.Pg.Port + "/" + app.conf.Pg.Db
	if app.db, err = pgsql.OpenConnection(pgsqlUrl); err != nil {
		log.Error("Connection Postgresql error ", err)
	}
	defer app.db.CloseConnection()

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.router = fastRouter.New()
	app.server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{
				Handler: nil,
			},
			Router: app.router,
		},
	}

	handlers.AddHandlers(app.router, app.db, &app.bots, app.server, &app.conf.Bot)

	go func() {
		if err = app.server.Start(app.conf.Bot.ListenAddress); err != nil {
			log.Error("Start server: \n", err)
		}
	}()

	app.bots = make(map[string]*bot.TBot)

	app.bots[botToken] = new(bot.TBot)
	app.bots[botToken].Bot, _ = bot.NewBot(&botToken)

	if err = app.bots[botToken].StartBot(app.conf.Bot.WebhookBase, app.conf.Bot.ListenAddress, app.server); err != nil {
		log.Error("Start bot\n", err)
	}

	// On exit, close, error program
	go func() {
		<-sigs
		log.Info("Stopping...")
		for _, v := range app.bots {
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
