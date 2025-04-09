package mybot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/resources"
	"log"
	"os"
)

type Bot struct {
	tgbotapi.BotAPI
}

func Register() *Bot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN environment variable not set")
	}
	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	_, err = tgbot.Request(tgbotapi.SetMyCommandsConfig{Commands: Commands})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", tgbot.Self.UserName)

	result := &Bot{BotAPI: *tgbot}

	StartReminderTask(result)

	return result
}

func (b *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	return b.BotAPI.GetUpdatesChan(u)
}

func (b *Bot) SendText(chatId int64, text string) *tgbotapi.Message {
	return b.sendWithOptions(chatId, text, res.NoPicture, nil, 0)
}

func (b *Bot) SendPic(chatId int64, text string, imageKey res.ImageKey) *tgbotapi.Message {
	return b.sendWithOptions(chatId, text, imageKey, nil, 0)
}

func (b *Bot) SendPicForceReply(chatId int64, text string, imageKey res.ImageKey, replyToMessageId int) *tgbotapi.Message {
	return b.sendWithOptions(chatId, text, imageKey, nil, replyToMessageId)
}

func (b *Bot) SendPicWithKeyboard(chatId int64, text string, imageKey res.ImageKey, keyboard *tgbotapi.InlineKeyboardMarkup) *tgbotapi.Message {
	return b.sendWithOptions(chatId, text, imageKey, keyboard, 0)
}

func (b *Bot) PinMessage(chatId int64, messageId int) {
	pinConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              chatId,
		MessageID:           messageId,
		DisableNotification: false,
	}

	_, err := b.BotAPI.Request(pinConfig)
	if err != nil {
		log.Println("[ERROR] Cannot pin message", err)
	}
}

//----------------------------------------------------------------------------------------

func (b *Bot) sendWithOptions(chatId int64, text string, imageKey res.ImageKey, keyboard *tgbotapi.InlineKeyboardMarkup, replyToMessageId int) *tgbotapi.Message {
	if keyboard != nil && replyToMessageId != 0 {
		panic("Одновременное задание keyboard и replyToMessageId запрещено!")
	}
	if file, ok := res.GetImage(imageKey); ok {
		return b.sendPicInternal(chatId, text, file, keyboard, replyToMessageId)
	} else {
		return b.sendTextInternal(chatId, text, keyboard, replyToMessageId)
	}
}

func (b *Bot) sendPicInternal(chatId int64, text string, file []byte, keyboard *tgbotapi.InlineKeyboardMarkup, replyToMessageId int) *tgbotapi.Message {
	photo := tgbotapi.FileBytes{
		Name:  "pic",
		Bytes: file,
	}

	message := tgbotapi.NewPhoto(chatId, photo)
	message.Caption = text
	message.ParseMode = tgbotapi.ModeMarkdown
	if keyboard != nil {
		message.ReplyMarkup = keyboard
	}
	if replyToMessageId != 0 {
		message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		message.ReplyToMessageID = replyToMessageId
	}

	return b.sendInternal(message)
}

func (b *Bot) sendTextInternal(chatId int64, text string, keyboard *tgbotapi.InlineKeyboardMarkup, replyToMessageId int) *tgbotapi.Message {
	//todo много дабликатов с отправкой картинки
	message := prepareMessage(chatId, text)
	if keyboard != nil {
		message.ReplyMarkup = keyboard
	}
	if replyToMessageId != 0 {
		message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		message.ReplyToMessageID = replyToMessageId
	}

	return b.sendInternal(message)
}

func prepareMessage(chatId int64, text string) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatId, text)
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
