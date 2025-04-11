package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

type RemindHandler struct {
}

const numberOfNames = 3

func (h RemindHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	users := sheets.Read()

	if len(users.AllUsers()) == 0 {
		msg := "Пока ещё никто не зарегистрировался 😢"

		bot.SendPic(chatID, msg, res.Sad)
	} else {

		nextBirthdayUsers, err := users.GetNextBirthdayUsers(numberOfNames)
		if err != nil {
			return err
		}

		maxNameLength := nextBirthdayUsers.GetMaxNameLength()
		maxMonthLength := nextBirthdayUsers.GetMaxMonthLength()

		log.Println("maxNameLength = ", maxNameLength)
		log.Println("maxMonthLength = ", maxMonthLength)

		msg := "Ближайшие именинники:\n```\n"
		for _, user := range nextBirthdayUsers.AllUsers() {
			msg += formatterStr(user, maxNameLength, maxMonthLength) + "\n"
		}
		msg += "```"
		bot.SendPic(chatID, msg, res.Many)
	}

	return nil
}

func formatterStr(user *usr.User, maxNameLength, maxMonthLength int) string {
	return user.FormattedString(maxNameLength, maxMonthLength) + formatDaysLeft(user.DaysBeforeBirthday())
}

func formatDaysLeft(days int) string {
	if days == 0 {
		return "(сегодня!😱)"
	} else {
		return fmt.Sprintf("(%d %s)", days, getDaysWord(days))
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

func (h RemindHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h RemindHandler) HandleCallback(*mybot.Bot, tgbotapi.Update, CallbackElement) error {
	return nil
}
