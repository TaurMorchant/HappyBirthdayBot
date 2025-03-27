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
	bot.Send(update.Message.Chat.ID, msg)

	return nil
}

// todo ограничить длиной самого длинного месяца
func formatterStr(user *usr.User, maxNameLength int) string {
	return user.FormattedString(maxNameLength) + formatDaysLeft(user.DaysBeforeBirthday())
}

func formatDaysLeft(days int) string {
	if days == 0 {
		return " (уже сегодня! 😱)"
	} else {
		return fmt.Sprintf(" (еще %d %s)", days, getDaysWord(days))
	}
}

func getDaysWord(n int) string {
	if n%100 >= 11 && n%100 <= 19 {
		return "дней"
	}
	switch n % 10 {
	case 1:
		return "день"
	case 2, 3, 4:
		return "дня"
	case 5, 6, 7, 8, 9, 0:
		return "дней"
	default:
		return ""
	}
}

func (h RemindHandler) HandleReply(*bot.Bot, tgbotapi.Update) error {
	return nil
}
