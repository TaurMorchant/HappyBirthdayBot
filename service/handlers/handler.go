package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/mybot"
)

type IHandler interface {
	Handle(bot *mybot.Bot, update tgbotapi.Update) error
	HandleReply(bot *mybot.Bot, update tgbotapi.Update) error
	HandleCallback(bot *mybot.Bot, update tgbotapi.Update, callback CallbackElement) error
}
