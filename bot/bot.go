package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	tgbotapi.BotAPI
}

// todo принимать чат и текст, добавлять markdown
func (b *Bot) SendWithEH(c tgbotapi.Chattable) tgbotapi.Message {
	mess, err := b.Send(c)
	if err != nil {
		log.Println("[ERROR] Cannot send message: ", err)
	}
	log.Printf("Message [%s] sent in chat %d", mess.Text, mess.Chat.ID)
	return mess
}
