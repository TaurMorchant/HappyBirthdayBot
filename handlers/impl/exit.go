package impl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
)

type ExitHandler struct {
}

func (h ExitHandler) Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()

	if _, ok := users.Get(usr.UserId(userID)); ok {
		users.Delete(usr.UserId(userID))
		sheets.Write(&users)
		bot.Send(tgbotapi.NewMessage(chatID, "Все пучком, ты удален из программы!"))
	} else {
		bot.Send(tgbotapi.NewMessage(chatID, "Слыш, ты и так не в программе!"))
	}
	return nil
}
