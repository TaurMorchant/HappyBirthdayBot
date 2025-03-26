package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
)

type StartHandler struct {
}

func (h StartHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	bot.SendWithEH(tgbotapi.NewMessage(chatID, "Я - бот-отеппибёздывватель! И вот что я умею:\n\nТУТ БУДЕТ ТЕКСТ"))

	return nil
}

func (h StartHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	return nil
}
