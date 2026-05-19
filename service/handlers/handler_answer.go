package handlers

import (
	"log"
	"strings"

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

	if len(update.Message.Photo) > 0 {
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		msg := tgbotapi.NewPhoto(mainChatId, tgbotapi.FileID(photo.FileID))
		msg.Caption = captionArguments(update.Message.Caption)
		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.BotAPI.Send(msg); err != nil {
			log.Println("[ERROR] Cannot forward photo:", err)
		}
		return nil
	}

	text := strings.NewReplacer("_", "\\_", "*", "\\*", "`", "\\`", "[", "\\[").Replace(update.Message.CommandArguments())
	bot.SendText(mainChatId, text)
	return nil
}

func captionArguments(caption string) string {
	idx := strings.Index(caption, " ")
	if idx == -1 {
		return ""
	}
	return caption[idx+1:]
}

func (h AnswerHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h AnswerHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
