package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/handlers"
	"happy-birthday-bot/handlers/impl"
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
		handleUpdate(bot, update)
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	defer handlePanic(bot, update)

	if update.Message == nil {
		return
	}

	log.Printf("Принято сообщение: %s", update.Message.Text)

	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == BotID {
		handleReply(bot, update)
	} else if update.Message.IsCommand() {
		handler, ok := handlers.Handlers[update.Message.Command()]
		if !ok {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Я не знаю команду '%s', откуда ты ее взял?", update.Message.Command())))
			return
		}
		err := handler.Handle(bot, update)
		if err != nil {
			message := fmt.Sprintf("Случилась какая-то неведомая фигня, напиши @morchant об этом, пожалуйста")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
		}
	}
}

func handlePanic(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	p := recover()
	err, _ := p.(error)
	log.Println("Panic: ", err)
	message := fmt.Sprintf("Случилась какая-то неведомая фигня, напиши @morchant об этом, пожалуйста")
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
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
