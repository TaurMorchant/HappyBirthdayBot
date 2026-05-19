package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/config"
	"happy-birthday-bot/mybot"
)

type RunReminderHandler struct{}

func (h RunReminderHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	if update.Message.Chat.ID != config.GetInt64Property(config.AdminChatIdProp) {
		return nil
	}
	mybot.RunReminderTask(bot)
	return nil
}

func (h RunReminderHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error { return nil }
func (h RunReminderHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
