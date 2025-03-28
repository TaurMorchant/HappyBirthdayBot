package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/handlers"
	"happy-birthday-bot/reminder"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/restrictions"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"
)

const BotID = 7947290853

func main() {
	file := configureLogger()
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Cannot close file!", err)
		}
	}(file)

	hapBirBot := registerBot()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := hapBirBot.GetUpdatesChan(u)

	reminder.StartReminderTask(hapBirBot)

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
		{Command: handlers.List, Description: "Все дни рождения"},
		{Command: handlers.Join, Description: "Присоединиться к программе"},
		{Command: handlers.Wishlist, Description: "Настроить свой Wishlist"},
		{Command: handlers.Reminders, Description: "Ближайшие дни рождения"},
		{Command: handlers.Exit, Description: "Выйти из программы"},
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

	if update.CallbackQuery != nil {
		log.Println("[TRACE] update.CallbackQuery.From.ID = ", update.CallbackQuery.From.ID)
		log.Println("[TRACE] update.CallbackQuery.Message.MessageID = ", update.CallbackQuery.Message.MessageID)
		log.Println("[TRACE] update.CallbackQuery.From.ID = ", update.CallbackQuery.From.ID)
		handleCallback(bot, update)
		return
	} else if update.Message != nil {
		log.Println("[TRACE] update.Message.From.ID = ", update.Message.From.ID)
		log.Println("[TRACE] update.Message.Chat.ID = ", update.Message.Chat.ID)

		if notRestricted(bot, update) {
			log.Printf("Принято сообщение: %s", update.Message.Text)

			if update.Message.IsCommand() {
				handler, ok := handlers.Handlers[update.Message.Command()]
				if !ok {
					bot.Send(update.Message.Chat.ID, fmt.Sprintf("Я не знаю команду '%s', откуда ты ее взял?", update.Message.Command()))
					return
				}
				err := handler.Handle(bot, update)
				if err != nil {
					panic(err)
				}
			} else if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == BotID {
				handleReply(bot, update)
			}
		}
	}
}

func notRestricted(bot *bot.Bot, update tgbotapi.Update) bool {
	if !restrictions.IsUserAllowed(update.Message.From.ID) {
		bot.Send(update.Message.Chat.ID, fmt.Sprintf("Прости %s, мне запрещено с тобой общаться!", update.Message.From.UserName))
		return false
	}

	if !restrictions.IsChatAllowed(update.Message.Chat.ID) {
		bot.Send(update.Message.Chat.ID, fmt.Sprintf("Прости, мне запрещено общаться в этом чате!"))
		return false
	}

	return true
}

func handlePanic(bot *bot.Bot, update tgbotapi.Update) {
	if p := recover(); p != nil {
		log.Println("[PANIC] Panic was catch: ", p)
		log.Println(string(debug.Stack()))
		message := fmt.Sprintf("Случилась какая-то неведомая фигня, напиши @morchant об этом, пожалуйста")
		bot.SendWithPicBasic(resolveChatId(update), message, res.Do_not_understand)
	}
}

func resolveChatId(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	} else {
		return update.Message.Chat.ID
	}
}

func handleReply(bot *bot.Bot, update tgbotapi.Update) {
	log.Println("Handle reply")

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	if handler, ok := handlers.WaitingForReplyHandlers.Get(userID); ok {
		err := handler.HandleReply(bot, update)
		if err != nil {
			log.Println("Error in reply:", err)

			bot.SendWithForceReply(chatID, messageID, err.Error())
			return
		}
		handlers.WaitingForReplyHandlers.Delete(userID)
	}
}

func handleCallback(bot *bot.Bot, update tgbotapi.Update) {
	log.Println("Handle callback")

	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	log.Println("userID = ", userID)

	removeButtonAnimation(bot, update)

	if callbackElement, ok := handlers.WaitingForCallbackHandlers.Get(messageId); ok {
		if callbackElement.UserId != userID {
			bot.SendWithPicBasic(chatId, "Это не для тебя кнопки, не трогай!", res.Angry_cats)
			return
		} else {
			err := callbackElement.Handler.HandleCallback(bot, update)
			if err != nil {
				log.Panic("Error in callback:", err)
				return
			}

			handlers.WaitingForCallbackHandlers.Delete(messageId)
		}
	}
	removeInlineButtons(bot, update)
}

func removeButtonAnimation(bot *bot.Bot, update tgbotapi.Update) {
	callbackCfg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := bot.Request(callbackCfg); err != nil {
		log.Panic("Ошибка при обработке callback:", err)
	}
}

func removeInlineButtons(bot *bot.Bot, update tgbotapi.Update) {
	editMsg := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)

	if _, err := bot.BotAPI.Send(editMsg); err != nil {
		log.Panic("Ошибка при редактировании сообщения:", err)
	}
}
