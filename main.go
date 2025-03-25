package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/handler"
	"happy-birthday-bot/handler/impl"
	"log"
	"os"
)

const BotID = 7947290853

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN environment variable not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // Игнорируем не сообщения
			log.Println("NOT MESSAGE ", update)
			continue
		}

		log.Printf("Принято сообщение: %s", update.Message.Text)

		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == BotID {
			handleReply(bot, update)
		} else if update.Message.IsCommand() {
			handler, ok := handler.Handlers[update.Message.Command()]
			if !ok {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Я не знаю команду '%s', откуда ты ее взял?", update.Message.Command())))
				continue
			}
			handler.Handle(bot, update)
		}
	}
}

func handleReply(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Println("Handle reply")

	chatID := update.Message.Chat.ID
	err := impl.HandleReply(bot, update)
	if err != nil {
		log.Println(err)

		message := tgbotapi.NewMessage(chatID, err.Error())
		message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		message.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(message)
	}
}
