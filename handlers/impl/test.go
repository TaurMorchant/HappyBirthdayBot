package impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TestHandler struct {
}

func (h TestHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	bot.Send(tgbotapi.NewMessage(chatID, "test"))

	return nil
}
