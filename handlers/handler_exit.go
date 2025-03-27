package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"strings"
)

type ExitHandler struct {
}

func (h ExitHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()

	if _, ok := users.Get(usr.UserId(userID)); ok {
		msg := "Ты точно уверен, что не хочешь быть отхеппибёзднутым?\n\nЕсли уверен, ответь на  это сообщение `Да`"
		bot.SendWithForceReply(chatID, msg)

		WaitForReply(usr.UserId(userID), h)
	} else {
		bot.Send(chatID, "Слыш, ты и так не в программе!")
	}
	return nil
}

func (h ExitHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if strings.EqualFold(update.Message.Text, "да") {
		users := sheets.Read()
		users.Delete(usr.UserId(userID))
		sheets.Write(&users)

		bot.Send(chatID, "Все пучком, ты удален из программы!")
	}
	return nil
}
