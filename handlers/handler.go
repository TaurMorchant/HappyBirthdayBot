package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/handlers/impl"
)

type IHandler interface {
	Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error
}

var Handlers = map[string]IHandler{
	Test:  &impl.TestHandler{},
	Start: &impl.StartHandler{},
	Join:  &impl.JoinHandler{},
	Exit:  &impl.ExitHandler{},
	List:  &impl.ListHandler{},
}
