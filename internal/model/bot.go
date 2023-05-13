package model

type BotStatus int

var (
	StatusBotActive BotStatus
	StatusBotRunnig BotStatus = 1
	StatusBotDel    BotStatus = 2
)

type Bot struct {
	Id     int64
	UserId int64
	Title  *string
	Token  *string
	Status BotStatus
}
