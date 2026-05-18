package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/config"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
)

type SelectHandler struct{}

func (h SelectHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	if update.Message.Chat.ID != config.GetInt64Property(config.AdminChatIdProp) {
		return nil
	}
	result, err := db.QueryRaw(update.Message.CommandArguments())
	if err != nil {
		bot.SendText(update.Message.Chat.ID, "```\nERROR: "+err.Error()+"\n```")
		return nil
	}
	bot.SendText(update.Message.Chat.ID, "```\n"+result+"\n```")
	return nil
}

func (h SelectHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error { return nil }
func (h SelectHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
