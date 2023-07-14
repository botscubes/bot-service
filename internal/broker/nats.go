package broker

import (
	"fmt"

	"github.com/botscubes/bot-service/internal/config"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

type NatsBroker struct {
	nc *nats.Conn
}

var (
	natsCodeOk = "200"
)

func NewNatsBroker(natsURL string) (*NatsBroker, error) {
	nc, err := nats.Connect(natsURL, nats.MaxReconnects(-1))
	if err != nil {
		return nil, err
	}

	return &NatsBroker{
		nc: nc,
	}, nil
}

func (b *NatsBroker) CloseConnection() {
	b.nc.Drain() //nolint:errcheck
}

type startBotPayload struct {
	BotId int64  `json:"botId"`
	Token string `json:"token"`
}

func (b *NatsBroker) StartBot(botId int64, token string) error {
	payload, err := json.Marshal(startBotPayload{
		BotId: botId,
		Token: token,
	})
	if err != nil {
		return err
	}

	res, err := b.nc.Request("worker.start", payload, config.NatsReqTimeout)
	if err != nil {
		return err
	}

	if string(res.Data) != natsCodeOk {
		return fmt.Errorf("nats get res error: %v", string(res.Data))
	}

	return nil
}

type stopBotPayload struct {
	BotId int64 `json:"botId"`
}

func (b *NatsBroker) StopBot(botId int64) error {
	payload, err := json.Marshal(stopBotPayload{
		BotId: botId,
	})
	if err != nil {
		return err
	}

	res, err := b.nc.Request("worker.stop", payload, config.NatsReqTimeout)
	if err != nil {
		return err
	}

	if string(res.Data) != natsCodeOk {
		return fmt.Errorf("nats get res error: %v", string(res.Data))
	}

	return nil
}
