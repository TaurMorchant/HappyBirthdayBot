package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"log"
)

type ListHandler struct {
}

func (h ListHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	log.Printf("Handle list command")
	chatID := update.Message.Chat.ID

	msg := "📅 Список всех участников:\n```\n"

	users := sheets.Read()
	usersSlice := users.GetAllUsers()

	maxNameLength := users.GetMaxNameLength()

	for _, user := range usersSlice {
		msg += user.FormattedString(maxNameLength) + "\n"
	}
	msg += "\n```"
	bot.SendWithPic(chatID, msg, res.Many_of_cats, nil)

	return nil
}

func (h ListHandler) HandleReply(*bot.Bot, tgbotapi.Update) error {
	return nil
}

func (h ListHandler) HandleCallback(*bot.Bot, tgbotapi.Update) error {
	return nil
}
