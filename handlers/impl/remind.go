package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"time"
)

type RemindHandler struct {
}

const numberOfNames = 3

func (h RemindHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	users := sheets.Read()
	//timeNow := time.Now()
	timeNow, _ := time.Parse("02.01.2006", "13.12.2025")

	var result [numberOfNames]string
	var i = 0

	maxNameLength := users.GetMaxNameLength()

	for _, user := range users.GetAllUsers() {
		if user.Birthday.Before(timeNow) {
			continue
		}
		if i < numberOfNames {
			result[i] = formatterStr(&user, timeNow, maxNameLength)
			i++
		} else {
			break
		}
	}
	lastIndex := i
	if lastIndex < numberOfNames-1 {
		timeNow = timeNow.AddDate(-1, 0, 0)
		for j := 0; j < numberOfNames-lastIndex; j++ {
			result[i] = formatterStr(&users.GetAllUsers()[j], timeNow, maxNameLength)
			i++
		}
	}

	msg := "Ближайшие именинники:\n```\n"
	for _, str := range result {
		msg += str + "\n"
	}
	msg += "```"
	message := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeMarkdown
	bot.SendWithEH(message)

	return nil
}

func formatterStr(user *usr.User, time time.Time, maxNameLength int) string {
	duration := user.Birthday.Sub(time)
	days := int(duration.Hours() / 24)
	return user.FormattedString(maxNameLength) + fmt.Sprintf(" (еще %d дней)", days)
}
