package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/handlers"
	"happy-birthday-bot/usr"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"
)

const BotID = 7947290853

func main() {
	file := configureLogger()
	defer file.Close()

	hapBirBot := registerBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := hapBirBot.GetUpdatesChan(u)

	for update := range updates {
		handleUpdate(hapBirBot, update)
	}
}

func configureLogger() *os.File {
	fileName := fmt.Sprintf("happy_birthday_bot_%s.log", time.Now().Format("2006-01-02_15.04.05"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	return file
}

func registerBot() *bot.Bot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN environment variable not set")
	}
	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	commands := []tgbotapi.BotCommand{
		{Command: handlers.Test, Description: "test"},
		{Command: handlers.List, Description: "Посмотреть всех в программе"},
		{Command: handlers.Join, Description: "Присоединиться к прогррамме"},
		{Command: handlers.Exit, Description: "Выйти из программы"},
		{Command: handlers.Reminders, Description: "Ближайшие дни рождения"},
	}

	_, err = tgbot.Request(tgbotapi.SetMyCommandsConfig{Commands: commands})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", tgbot.Self.UserName)

	return &bot.Bot{BotAPI: *tgbot}
}

func handleUpdate(bot *bot.Bot, update tgbotapi.Update) {
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
			bot.SendWithEH(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Я не знаю команду '%s', откуда ты ее взял?", update.Message.Command())))
			return
		}
		err := handler.Handle(bot, update)
		if err != nil {
			message := fmt.Sprintf("Случилась какая-то неведомая фигня, напиши @morchant об этом, пожалуйста")
			bot.SendWithEH(tgbotapi.NewMessage(update.Message.Chat.ID, message))
		}
	}
}

func handlePanic(bot *bot.Bot, update tgbotapi.Update) {
	p := recover()
	if p == nil {
		return
	}
	if err, ok := p.(error); ok {
		log.Println("Panic: ", err)
		fmt.Println(string(debug.Stack()))
		message := fmt.Sprintf("Случилась какая-то неведомая фигня, напиши @morchant об этом, пожалуйста")
		bot.SendWithEH(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	}
}

func handleReply(bot *bot.Bot, update tgbotapi.Update) {
	log.Println("Handle reply")

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	handler := handlers.GetWaitingHandler(usr.UserId(userID))
	if handler == nil {
		return
	} else {
		err := handler.HandleReply(bot, update)
		if err != nil {
			log.Println(err)

			message := tgbotapi.NewMessage(chatID, err.Error())
			message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
			message.ParseMode = tgbotapi.ModeMarkdown
			bot.SendWithEH(message)
			return
		}
		handlers.RemoveWaitingHandler(usr.UserId(userID))
	}

	//
	//err := handlers2.HandleReply(bot, update)
	//if err != nil {
	//	log.Println(err)
	//
	//	message := tgbotapi.NewMessage(chatID, err.Error())
	//	message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	//	message.ParseMode = tgbotapi.ModeMarkdown
	//	bot.SendWithEH(message)
	//}
}
