package impl

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/date"
	"happy-birthday-bot/sheets"
)

type ListHandler struct {
}

func (h ListHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	msg := "📅 Список всех участников:\n```\n"

	users := sheets.Read()
	usersSlice := users.GetAllUsers()

	for _, user := range usersSlice {
		msg += fmt.Sprintf("🎂 %*s — %10s\n", users.GetMaxNameLength(), user.Name, date.FormatDate(user.Birthday))
	}
	msg += "\n```"
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "markdown"
	bot.Send(message)

	return nil
}
