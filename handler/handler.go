package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/handler/impl"
)

type IHandler interface {
	Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

var Handlers = map[string]IHandler{
	Test: &impl.TestHandler{},
	Join: &impl.JoinHandler{},
	Exit: &impl.ExitHandler{},
	List: &impl.ListHandler{},
}
