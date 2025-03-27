package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
)

type RemindHandler struct {
}

const numberOfNames = 3

func (h RemindHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	users := sheets.Read()

	maxNameLength := users.GetMaxNameLength()

	nextBirthdayUsers, err := users.GetNextBirthdayUsers(numberOfNames)
	if err != nil {
		return err
	}

	msg := "Ближайшие именинники:\n```\n"
	for _, user := range nextBirthdayUsers {
		msg += formatterStr(user, maxNameLength) + "\n"
	}
	msg += "```"
	message := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeMarkdown
	bot.SendWithEH(message)

	return nil
}

// todo ограничить длиной самого длинного месяца
func formatterStr(user *usr.User, maxNameLength int) string {
	days := user.DaysBeforeBirthday()
	return user.FormattedString(maxNameLength) + fmt.Sprintf(" (еще %d дней)", days)
}

func (h RemindHandler) HandleReply(*bot.Bot, tgbotapi.Update) error {
	return nil
}
