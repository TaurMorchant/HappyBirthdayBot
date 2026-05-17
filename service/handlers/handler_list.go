package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/db"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"log"
)

type ListHandler struct {
}

func (h ListHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	log.Printf("Handle list command")
	chatID := update.Message.Chat.ID

	users := db.ReadUsers()
	usersSlice := users.AllUsers()

	if len(usersSlice) == 0 {
		msg := "Пока ещё никто не зарегистрировался 😢"

		bot.SendPic(chatID, msg, res.Sad)
	} else {
		msg := "📅 Вот список всех участников:\n```\n"

		maxNameLength := users.GetMaxNameLength()
		maxMonthLength := users.GetMaxMonthLength()

		for _, user := range usersSlice {
			msg += user.FormattedString(maxNameLength, maxMonthLength) + "\n"
		}
		msg += "\n```"
		bot.SendPic(chatID, msg, res.Many)
	}

	return nil
}

func (h ListHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h ListHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
