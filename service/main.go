package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/config"
	"happy-birthday-bot/handlers"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
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

	config.InitConfigs()
	sheets.InitSpreadsheetService()

	bot := mybot.Register()

	for update := range bot.GetUpdatesChan() {
		handleUpdate(bot, update)
	}
}

func configureLogger() *os.File {
	fileName := fmt.Sprintf("happy_birthday_bot_%s.log", time.Now().Format("2006-01-02_15.04.05"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	return file
}

func handleUpdate(bot *mybot.Bot, update tgbotapi.Update) {
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
		log.Println("[TRACE] update.Message.Text = ", update.Message.Text)

		if notRestricted(bot, update) {
			if update.Message.IsCommand() {
				handler, ok := handlers.Handlers[update.Message.Command()]
				if !ok {
					msg := fmt.Sprintf("Я не знаю команду '%s', откуда ты ее взял?", update.Message.Command())
					bot.SendPic(update.Message.Chat.ID, msg, res.Error)
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

func notRestricted(bot *mybot.Bot, update tgbotapi.Update) bool {
	if !config.IsUserAllowed(update.Message.From.ID) {
		msg := fmt.Sprintf("Прости %s, мне запрещено с тобой общаться!", update.Message.From.UserName)
		bot.SendPic(update.Message.Chat.ID, msg, res.Error)
		return false
	}

	if !config.IsChatAllowed(update.Message.Chat.ID) {
		log.Printf("Chat %d is not allowed!", update.Message.Chat.ID)
		msg := fmt.Sprintf("Прости, мне запрещено общаться в этом чате!")
		bot.SendPic(update.Message.Chat.ID, msg, res.Error)
		return false
	}

	return true
}

func handlePanic(bot *mybot.Bot, update tgbotapi.Update) {
	defer handleFinalPanic()
	if p := recover(); p != nil {
		log.Println("[PANIC] Panic was catch: ", p)
		log.Println(string(debug.Stack()))
		message := fmt.Sprintf("Случилась какая-то неведомая фигня! 😱\n\n" +
			"Напиши @morchant об этом, пожалуйста")
		bot.SendPic(resolveChatId(update), message, res.Error)
	}
}

func handleFinalPanic() {
	if p := recover(); p != nil {
		log.Println("[PANIC] Final Panic was catch: ", p)
		log.Println(string(debug.Stack()))
	}
}

func resolveChatId(update tgbotapi.Update) int64 {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	} else {
		return update.Message.Chat.ID
	}
}

func handleReply(bot *mybot.Bot, update tgbotapi.Update) {
	log.Println("Handle reply")

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	if handler, ok := handlers.WaitingForReplyHandlers.Get(userID); ok {
		err := handler.HandleReply(bot, update)
		if err != nil {
			log.Println("Error in reply:", err)

			bot.SendPicForceReply(chatID, err.Error(), res.Error, messageID)
			return
		}
		handlers.WaitingForReplyHandlers.Delete(userID)
	}
}

func handleCallback(bot *mybot.Bot, update tgbotapi.Update) {
	log.Println("Handle callback")

	chatId := update.CallbackQuery.Message.Chat.ID
	messageId := update.CallbackQuery.Message.MessageID
	userID := update.CallbackQuery.From.ID

	log.Println("userID = ", userID)

	removeButtonAnimation(bot, update)

	if callbackElement, ok := handlers.WaitingForCallbackHandlers.Get(messageId); ok {
		if callbackElement.UserId != userID {
			bot.SendPic(chatId, "Это не для тебя кнопки, не трогай!", res.Angry)
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

func removeButtonAnimation(bot *mybot.Bot, update tgbotapi.Update) {
	callbackCfg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := bot.Request(callbackCfg); err != nil {
		log.Panic("Ошибка при обработке callback:", err)
	}
}

func removeInlineButtons(bot *mybot.Bot, update tgbotapi.Update) {
	editMsg := tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
	)

	if _, err := bot.BotAPI.Send(editMsg); err != nil {
		log.Panic("Ошибка при редактировании сообщения:", err)
	}
}
