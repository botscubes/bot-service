package model

// Statuses
type BotStatus int

var (
	StatusBotActive BotStatus = 0
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
