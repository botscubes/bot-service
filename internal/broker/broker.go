package broker

type Broker interface {
	StartBot(botId int64, token string) error
	StopBot(botId int64) error
	CloseConnection()
}
