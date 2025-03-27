package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

type ExitHandler struct {
}

const okButton = "ok"
const cancelButton = "cancel"

func (h ExitHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()

	if _, ok := users.Get(usr.UserId(userID)); ok {
		msg := "Ты точно уверен, что хочешь уйти из программы и не жулаешь быть отхеппибёзднутым?"

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да!", okButton),
				tgbotapi.NewInlineKeyboardButtonData("ГАЛЯ, ОТМЕНА!!!", cancelButton),
			),
		)

		sentMessage := bot.SendWithKeyboard(chatID, msg, inlineKeyboard)

		//bot.SendWithForceReply(chatID, msg)

		WaitForCallback(sentMessage.MessageID, userID, h)
	} else {
		bot.Send(chatID, "Слыш, ты и так не в программе!")
	}
	return nil
}

func (h ExitHandler) HandleReply(*bot.Bot, tgbotapi.Update) error {
	return nil
}

func (h ExitHandler) HandleCallback(bot *bot.Bot, update tgbotapi.Update) error {
	log.Println("Handle callback for ExitHandler")
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	if update.CallbackQuery.Data == okButton {
		users := sheets.Read()
		users.Delete(usr.UserId(userID))
		sheets.Write(&users)

		bot.Send(chatID, "Все пучком, ты удален из программы!")
	} else if update.CallbackQuery.Data == cancelButton {
		bot.Send(chatID, "Да ладно, ладно, не ори!")
	} else {
		bot.Send(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!")
	}

	return nil
}
