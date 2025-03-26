package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
)

type TestHandler struct {
}

func (h TestHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	bot.SendWithEH(tgbotapi.NewMessage(chatID, "test"))

	return nil
}

func (h TestHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	return nil
}
