package impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
}

func (h StartHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	bot.Send(tgbotapi.NewMessage(chatID, "Я - бот-отеппибёздывватель! И вот что я умею:\n\nТУТ БУДЕТ ТЕКСТ"))

	return nil
}
