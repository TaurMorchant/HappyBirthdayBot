package impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TestHandler struct {
}

func (h TestHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	log.Printf("Message INFO: [%s:%s] %s", update.Message.From.UserName, userID, update.Message.Text)

	bot.Send(tgbotapi.NewMessage(chatID, "test"))
}
