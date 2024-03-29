package bot

import (
	"strconv"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/mymmrac/telego"
	"go.uber.org/zap"
)

type BotService struct {
	conf *config.ServiceConfig
	log  *zap.SugaredLogger
}

func NewBotService(c *config.ServiceConfig, l *zap.SugaredLogger) *BotService {
	return &BotService{
		conf: c,
		log:  l,
	}
}

func (bs *BotService) StartBot(botId int64, token string) error {
	bot, err := telego.NewBot(token, telego.WithHealthCheck())
	if err != nil {
		bs.log.Errorw("failed telego newBot", "error", err)
		return err
	}

	// Add remove already exists webhook

	return bot.SetWebhook(&telego.SetWebhookParams{
		URL: "https://" + bs.conf.WebhookDomain + bs.conf.WebhookPath + strconv.FormatInt(botId, 10),
	})
}

func (bs *BotService) StopBot(token string) error {
	bot, err := telego.NewBot(token, telego.WithHealthCheck())
	if err != nil {
		bs.log.Errorw("failed telego newBot", "error", err)
		return err
	}

	return bot.DeleteWebhook(nil)
}

func (bs *BotService) TokenHealthCheck(token string) (bool, error) {
	_, err := telego.NewBot(token, telego.WithHealthCheck())
	if err != nil {
		if err.Error() == ErrTgAuth401.Error() {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
