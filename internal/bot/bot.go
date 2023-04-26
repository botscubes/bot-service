package bot

import (
	"time"

	"github.com/botscubes/bot-service/pkg/log"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TBot struct {
	Bot     *telego.Bot
	Updates <-chan telego.Update
	Handler *th.BotHandler
}

func NewBot(token *string) (*telego.Bot, error) {
	return telego.NewBot(*token, telego.WithHealthCheck(), telego.WithDefaultDebugLogger())
}

func (btx *TBot) setBotHandler() {
	btx.Handler.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		var err error = nil
		chatID := tu.ID(message.Chat.ID)
		_, err = bot.CopyMessage(tu.CopyMessage(chatID, chatID, message.MessageID))
		if err != nil {
			bot.Logger().Errorf("Failed to copy message: %s", err)
		}

		bot.Logger().Debugf("Copied message with ID %d in chat %d", message.MessageID, chatID.ID)
	})
}

func (btx *TBot) startBotHandler() {
	go btx.Handler.Start()
}

func (btx *TBot) StartBot(webhookBase string, listenAddress string, server *telego.MultiBotWebhookServer) error {
	var err error = nil

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
		if btx.Handler, err = th.NewBotHandler(btx.Bot, btx.Updates, th.WithStopTimeout(time.Second*10)); err != nil {
			return err
		}

		btx.setBotHandler()
	}

	btx.startBotHandler()

	if !btx.Bot.IsRunningWebhook() {
		go func(b *telego.Bot) {
			if err := b.StartWebhook(listenAddress); err != nil {
				log.Error("Start webhook:", err)
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

	if err := btx.Bot.DeleteWebhook(nil); err != nil {
		return err
	}

	return nil
}
