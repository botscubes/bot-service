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
		URL: webhookBase + "/bot" + btx.Bot.Token(),
	})

	btx.Updates, err = btx.Bot.UpdatesViaWebhook(
		"/bot"+btx.Bot.Token(),
		telego.WithWebhookServer(server),
	)
	if err != nil {
		return err
	}

	btx.Handler, err = th.NewBotHandler(btx.Bot, btx.Updates, th.WithStopTimeout(time.Second*10))
	if err != nil {
		return err
	}

	btx.setBotHandler()

	btx.startBotHandler()

	go func(b *telego.Bot) {
		err := b.StartWebhook(listenAddress)
		if err != nil {
			log.Error("Start webhook:", err)
		}
	}(btx.Bot)

	return nil
}
