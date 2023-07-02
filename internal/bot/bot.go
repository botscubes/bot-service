package bot

import (
	"time"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

type BotService struct {
	bots   map[int64]*TBot
	conf   *config.BotConfig
	server *telego.MultiBotWebhookServer
}

type TBot struct {
	Id        int64
	Rdb       *rdb.Rdb
	Db        *pgsql.Db
	Bot       *telego.Bot
	Updates   <-chan telego.Update
	Handler   *th.BotHandler
	log       *zap.SugaredLogger
	IsRunning bool
}

const handlerTimeout = 10 // sec

func NewBotService(c *config.BotConfig, s *telego.MultiBotWebhookServer) *BotService {
	return &BotService{
		conf:   c,
		server: s,
		bots:   make(map[int64]*TBot),
	}
}

func (bs *BotService) NewBot(
	token *string,
	botId int64,
	logger *zap.SugaredLogger,
	r *rdb.Rdb,
	db *pgsql.Db,
) error {
	bot, err := telego.NewBot(*token, telego.WithHealthCheck(), telego.WithDefaultDebugLogger())
	if err != nil {
		return err
	}

	bs.bots[botId] = &TBot{
		Id:        botId,
		Bot:       bot,
		log:       logger,
		Rdb:       r,
		Db:        db,
		IsRunning: false,
	}

	return nil
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

func (btx *TBot) startBot(webhookBase string, listenAddress string, server *telego.MultiBotWebhookServer) error {
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

	btx.IsRunning = true

	return nil
}

func (btx *TBot) stopBot(stopWebhookServer bool) error {
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

	btx.IsRunning = false

	return btx.Bot.DeleteWebhook(nil)
}

func (bs *BotService) StopBots() error {
	for _, v := range bs.bots {
		if err := v.stopBot(true); err != nil {
			return err
		}
	}

	return nil
}

// TODO: -> check runnig
func (bs *BotService) CheckBotExist(botID int64) bool {
	_, ok := bs.bots[botID]
	return ok
}

func (bs *BotService) StartBot(botID int64) error {
	bot, ok := bs.bots[botID]
	if !ok {
		return ErrBotNotFound
	}

	bot.log.Debug("Call Bot Start Handler")

	// return bot.startBot(bs.conf.WebhookBase, bs.conf.ListenAddress, bs.server)

	return nil
}

func (bs *BotService) StopBot(botID int64) error {
	bot, ok := bs.bots[botID]
	if !ok {
		return ErrBotNotFound
	}

	return bot.stopBot(false)
}

func (bs *BotService) BotIsRunnig(botID int64) (bool, error) {
	if ok := bs.CheckBotExist(botID); !ok {
		return false, ErrBotNotFound
	}

	return bs.bots[botID].IsRunning, nil
}
