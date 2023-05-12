package bot

import (
	"time"

	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

type TBot struct {
	Id      int64
	Rdb     *rdb.Rdb
	Db      *pgsql.Db
	Bot     *telego.Bot
	Updates <-chan telego.Update
	Handler *th.BotHandler
	log     *zap.SugaredLogger
}

const handlerTimeout = 10 // sec

func New(token *string, botId int64, logger *zap.SugaredLogger) (*TBot, error) {
	bot, err := telego.NewBot(*token, telego.WithHealthCheck(), telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}
	res := new(TBot)
	res.Id = botId
	res.Bot = bot
	res.log = logger
	return res, nil
}

func (btx *TBot) setMiddlwares() {
	btx.Handler.Use(th.PanicRecovery)
	btx.Handler.Use(btx.regUserMW)
}

func (btx *TBot) setHandlers() {
	// Handle command
	btx.Handler.Handle(btx.commandHandler(),
		th.Union(
			th.AnyCommand(),
		))

	// Handle message
	btx.Handler.Handle(btx.messageHandler(),
		th.Union(
			th.AnyMessage(),
			th.AnyEditedMessage(),
		))
}

func (btx *TBot) startHandler() {
	go btx.Handler.Start()
}

func (btx *TBot) StartBot(webhookBase string, listenAddress string, server *telego.MultiBotWebhookServer) error {
	var err error

	_ = btx.Bot.SetWebhook(&telego.SetWebhookParams{
		URL: webhookBase + "/webhook/bot" + btx.Bot.Token(),
	})

	if btx.Updates == nil {
		if btx.Updates, err = btx.Bot.UpdatesViaWebhook(
			"/webhook/bot"+btx.Bot.Token(),
			telego.WithWebhookServer(server),
		); err != nil {
			return err
		}
	}

	if btx.Handler == nil {
		if btx.Handler, err = th.NewBotHandler(btx.Bot, btx.Updates, th.WithStopTimeout(time.Second*handlerTimeout)); err != nil {
			return err
		}

		btx.setMiddlwares()
		btx.setHandlers()
	}

	btx.startHandler()

	if !btx.Bot.IsRunningWebhook() {
		go func(b *telego.Bot) {
			if err := b.StartWebhook(listenAddress); err != nil {
				btx.log.Error("Start webhook:", err)
			}
		}(btx.Bot)
	}

	return nil
}

func (btx *TBot) StopBot(stopWebhookServer bool) error {
	// WARN: The bot is not removed from the app.bots.
	// Because btx.Updates and btx.Handler
	// cannot be initialized again when StartBot is called again.
	// The reason is that fasthttp/router does not allow you to delete handlers.
	//
	// TODO: set state bot stopped (After call this func).
	// !!!! Find a way to delete handlers (UpdatesViaWebhook) !!!!

	if stopWebhookServer {
		if err := btx.Bot.StopWebhook(); err != nil {
			return err
		}
	}

	if btx.Handler != nil {
		btx.Handler.Stop()
	}

	return btx.Bot.DeleteWebhook(nil)
}

func (btx *TBot) SetDb(db *pgsql.Db) {
	btx.Db = db
}

func (btx *TBot) SetRdb(r *rdb.Rdb) {
	btx.Rdb = r
}
