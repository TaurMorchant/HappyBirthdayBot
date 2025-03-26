package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	tgbotapi.BotAPI
}

func (b *Bot) SendWithEH(c tgbotapi.Chattable) tgbotapi.Message {
	mess, err := b.Send(c)
	if err != nil {
		log.Println("[ERROR] Cannot send message: ", err)
	}
	return mess
}
