package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/valyala/fasthttp"

	fastRouter "github.com/fasthttp/router"

	"github.com/botscubes/bot-service/internal/bot"
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
}

var app App

const envPrefix = "TBOT_"

func Run() {

	var err error = nil

	webhookBase, ok := env("WEBHOOK_BASE")
	assert(ok, "Environment variable "+envPrefix+"WEBHOOK_BASE not found")

	listenAddress, ok := env("LISTEN_ADDRESS")
	assert(ok, "Environment variable "+envPrefix+"LISTEN_ADDRESS not found")

	botToken, ok := env("TOKEN")
	assert(ok, "Environment variable "+envPrefix+"TOKEN not found")

	app.config = new(appConfig)
	app.config.webhookBase = webhookBase
	app.config.listenAddress = listenAddress

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.router = createRouter()
	initHandlers()

	app.server = &telego.MultiBotWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Server: &fasthttp.Server{},
			Router: app.router,
		},
	}

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
	return os.LookupEnv(envPrefix + name)
}

// Check ok and exit program if ok is false
func assert(ok bool, args ...any) {
	if !ok {
		log.Fatal(append([]any{"FATAL:"}, args...)...)
	}
}
