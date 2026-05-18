package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/config"
	"happy-birthday-bot/mybot"
)

type AnswerHandler struct{}

func (h AnswerHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	if update.Message.Chat.ID != config.GetInt64Property(config.AdminChatIdProp) {
		return nil
	}
	mainChatId := config.GetInt64Property(config.MainChatIdProp)
	bot.SendText(mainChatId, update.Message.CommandArguments())
	return nil
}

func (h AnswerHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h AnswerHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
