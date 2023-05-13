package model

type BotStatus int

var (
	StatusBotActive BotStatus
)

type Bot struct {
	Id     int64
	UserId int64
	Title  *string
	Token  *string
	Status BotStatus
}
