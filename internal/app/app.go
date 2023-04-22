package app

// net package NOT USED
import (
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"

	fastRouter "github.com/fasthttp/router"

	"github.com/botscubes/bot-service/internal/api/handlers"
	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	"github.com/botscubes/bot-service/internal/net/client"
	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
)

type appConfig struct {
	webhookBase   string
	listenAddress string
}

type App struct {
	router *fastRouter.Router
	server *telego.MultiBotWebhookServer
	bots   map[string]*bot.TBot
	config *appConfig
	client *client.TClient
	db     *pgsql.Db
}

var app App

const envPrefix = "TBOT_"

func Run() {

	log.Debug("Init")

	var err error = nil

	// TODO: create util for get env, for example: getenv(name, prefix, check required)
	// MB Switch to yml configs

	// TODO: create server 50x error response

	webhookBase, ok := env(envPrefix + "WEBHOOK_BASE")
	assert(ok, "Environment variable "+envPrefix+"WEBHOOK_BASE not found")

	listenAddress, ok := env(envPrefix + "LISTEN_ADDRESS")
	assert(ok, "Environment variable "+envPrefix+"LISTEN_ADDRESS not found")

	botToken, ok := env(envPrefix + "TOKEN")
	assert(ok, "Environment variable "+envPrefix+"TOKEN not found")

	pgsqlBase, ok := env("POSTGRES_DB")
	assert(ok, "Environment variable POSTGRES_DB not found")

	pgsqlUser, ok := env("POSTGRES_USER")
	assert(ok, "Environment variable POSTGRES_USER not found")

	pgsqlPass, ok := env("POSTGRES_PASSWORD")
	assert(ok, "Environment variable POSTGRES_PASSWORD not found")

	pgsqlHost, ok := env("POSTGRES_HOST")
	assert(ok, "Environment variable POSTGRES_HOST not found")

	pgsqlPort, ok := env("POSTGRES_PORT")
	assert(ok, "Environment variable POSTGRES_PORT not found")

	pgsqlUrl := "postgres://" + pgsqlUser + ":" + pgsqlPass + "@" + pgsqlHost + ":" + pgsqlPort + "/" + pgsqlBase

	app.db, err = pgsql.OpenConnection(pgsqlUrl)
	if err != nil {
		log.Error("Connection Postgresql error ", err)
	}

	defer app.db.CloseConnection()

	app.db.GetTest()

	app.config = new(appConfig)
	app.config.webhookBase = webhookBase
	app.config.listenAddress = listenAddress

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.router = fastRouter.New()
	handlers.AddHandlers(app.router, app.db)

	app.server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{},
			Router: app.router,
		},
	}

	go func() {
		err = app.server.Start(listenAddress)
		if err != nil {
			log.Error("Start server ", err)
		}
	}()

	// UNUSED !!
	app.client = client.NewClient()

	app.bots = make(map[string]*bot.TBot)

	app.bots[botToken] = new(bot.TBot)
	app.bots[botToken].Bot, _ = bot.NewBot(botToken)

	err = app.bots[botToken].StartBot(app.config.webhookBase, app.config.listenAddress, app.server)
	if err != nil {
		log.Error("Start bot ", err)
	}

	// On exit, close, error program
	go func() {
		<-sigs
		log.Info("Stopping...")
		var err error = nil
		for _, v := range app.bots {
			err = v.Bot.StopWebhook()
			if err != nil {
				log.Error("Stop webhook:", err)
			}

			if v.Handler != nil {
				v.Handler.Stop()
			}

			err = v.Bot.DeleteWebhook(nil)
			if err != nil {
				log.Error("Delete webhook:", err)
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
