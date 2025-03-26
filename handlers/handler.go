package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	handlers "happy-birthday-bot/handlers/impl"
)

type IHandler interface {
	Handle(bot *bot.Bot, update tgbotapi.Update) error
}

var Handlers = map[string]IHandler{
	Test:  &handlers.TestHandler{},
	Start: &handlers.StartHandler{},
	Join:  &handlers.JoinHandler{},
	Exit:  &handlers.ExitHandler{},
	List:  &handlers.ListHandler{},
}
