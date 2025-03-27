package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	tgbotapi.BotAPI
}

func (b *Bot) Send(chatId int64, str string) *tgbotapi.Message {
	message := prepareMessage(chatId, str)
	return b.sendInternal(message)
}

func (b *Bot) SendWithForceReply(chatId int64, str string) *tgbotapi.Message {
	message := prepareMessage(chatId, str)
	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	return b.sendInternal(message)
}

func (b *Bot) SendWithKeyboard(chatId int64, str string, keyboard tgbotapi.InlineKeyboardMarkup) *tgbotapi.Message {
	message := prepareMessage(chatId, str)
	message.ReplyMarkup = keyboard
	return b.sendInternal(message)
}

//----------------------------------------------------------------------------------------

func prepareMessage(chatId int64, str string) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatId, str)
	message.ParseMode = tgbotapi.ModeMarkdown
	return &message
}

func (b *Bot) sendInternal(message tgbotapi.Chattable) *tgbotapi.Message {
	mess, err := b.BotAPI.Send(message)
	if err != nil {
		log.Panicln("[ERROR] Cannot send message: ", err)
		return nil
	}
	log.Printf("Message [%s] sent in chat %d", mess.Text, mess.Chat.ID)
	return &mess
}
