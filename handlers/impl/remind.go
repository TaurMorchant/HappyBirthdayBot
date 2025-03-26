package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
)

type RemindHandler struct {
}

func (h RemindHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {

	return nil
}
