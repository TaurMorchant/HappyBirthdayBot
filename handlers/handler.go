package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
)

type IHandler interface {
	Handle(bot *bot.Bot, update tgbotapi.Update) error
	HandleReply(bot *bot.Bot, update tgbotapi.Update) error
	HandleCallback(bot *bot.Bot, update tgbotapi.Update) error
}
