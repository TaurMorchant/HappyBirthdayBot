package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/resources"
	"log"
)

type Bot struct {
	tgbotapi.BotAPI
}

func (b *Bot) Send(chatId int64, msg string) *tgbotapi.Message {
	message := prepareMessage(chatId, msg)
	return b.sendInternal(message)
}

func (b *Bot) SendWithForceReply(chatId int64, msg string) *tgbotapi.Message {
	message := prepareMessage(chatId, msg)
	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	return b.sendInternal(message)
}

func (b *Bot) SendWithPic(chatId int64, msg string, imageKey res.ImageKey, keyboard *tgbotapi.InlineKeyboardMarkup, forceReply bool) *tgbotapi.Message {
	file, err := res.GetImage(imageKey)
	if err != nil {
		log.Println("Cannot get image", err)
		return b.sendWithKeyboard(chatId, msg, keyboard)
	} else {
		photo := tgbotapi.FileBytes{
			Name:  string(imageKey),
			Bytes: file,
		}

		message := tgbotapi.NewPhoto(chatId, photo)
		message.Caption = msg
		message.ParseMode = tgbotapi.ModeMarkdown
		if keyboard != nil {
			message.ReplyMarkup = keyboard
		}
		if forceReply {
			message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		}

		return b.sendInternal(message)
	}
}

//----------------------------------------------------------------------------------------

func (b *Bot) sendWithKeyboard(chatId int64, msg string, keyboard *tgbotapi.InlineKeyboardMarkup) *tgbotapi.Message {
	message := prepareMessage(chatId, msg)
	if keyboard != nil {
		message.ReplyMarkup = keyboard
	}
	return b.sendInternal(message)
}

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
