package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/config"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
)

type ExecHandler struct{}

func (h ExecHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	if update.Message.Chat.ID != config.GetInt64Property(config.AdminChatIdProp) {
		return nil
	}
	result, err := db.ExecRaw(update.Message.CommandArguments())
	if err != nil {
		bot.SendText(update.Message.Chat.ID, "```\nERROR: "+err.Error()+"\n```")
		return nil
	}
	bot.SendText(update.Message.Chat.ID, "```\n"+result+"\n```")
	return nil
}

func (h ExecHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error { return nil }
func (h ExecHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
