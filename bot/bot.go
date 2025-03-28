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

func (b *Bot) SendWithForceReply(chatId int64, replyToMessageId int, msg string) *tgbotapi.Message {
	message := prepareMessage(chatId, msg)
	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	message.ReplyToMessageID = replyToMessageId
	return b.sendInternal(message)
}

func (b *Bot) SendWithPicBasic(chatId int64, msg string, imageKey res.ImageKey) *tgbotapi.Message {
	return b.sendWithPicInternal(chatId, msg, imageKey, nil, 0)
}

func (b *Bot) SendWithPicAndForceReply(chatId int64, msg string, imageKey res.ImageKey, replyToMessageId int) *tgbotapi.Message {
	return b.sendWithPicInternal(chatId, msg, imageKey, nil, replyToMessageId)
}

func (b *Bot) SendWithPicAndKeyboard(chatId int64, msg string, imageKey res.ImageKey, keyboard *tgbotapi.InlineKeyboardMarkup) *tgbotapi.Message {
	return b.sendWithPicInternal(chatId, msg, imageKey, keyboard, 0)
}

//----------------------------------------------------------------------------------------

func (b *Bot) sendWithPicInternal(chatId int64, msg string, imageKey res.ImageKey, keyboard *tgbotapi.InlineKeyboardMarkup, replyToMessageId int) *tgbotapi.Message {
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
		if keyboard != nil && replyToMessageId != 0 {
			panic("Одновременное задание keyboard и replyToMessageId запрещено!")
		}
		if keyboard != nil {
			message.ReplyMarkup = keyboard
		}
		if replyToMessageId != 0 {
			message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
			message.ReplyToMessageID = replyToMessageId
		}

		return b.sendInternal(message)
	}
}

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
