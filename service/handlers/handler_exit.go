package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/mybot"
	res "happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

type ExitHandler struct {
}

const okButton = "ok"
const cancelButton = "cancel"

func (h ExitHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageId := update.Message.MessageID

	users := sheets.Read()

	if _, ok := users.Get(usr.UserId(userID)); ok {
		msg := "Ты уверен, что хочешь уйти из программы и не желаешь быть отхеппибёзднутым?"

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да!", okButton),
				tgbotapi.NewInlineKeyboardButtonData("ГАЛЯ, ОТМЕНА!!!", cancelButton),
			),
		)

		sentMessage := bot.SendPicWithKeyboard(chatID, msg, res.Sad, &inlineKeyboard, messageId)

		WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: h})
	} else {
		bot.SendPic(chatID, "Слыш, ты и так не в программе!", res.Suspicious)
	}
	return nil
}

func (h ExitHandler) HandleReply(*mybot.Bot, tgbotapi.Update) error {
	return nil
}

func (h ExitHandler) HandleCallback(bot *mybot.Bot, update tgbotapi.Update, _ CallbackElement) error {
	log.Println("Handle callback for ExitHandler")
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID

	if update.CallbackQuery.Data == okButton {
		users := sheets.Read()
		users.Delete(usr.UserId(userID))
		sheets.Write(&users)

		bot.SendPic(chatID, "Штош, ты удален", res.Sad)
	} else if update.CallbackQuery.Data == cancelButton {
		bot.SendPic(chatID, "Да ладно, ладно, не ори!", res.DoNotScream)
	} else {
		bot.SendPic(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!", res.Error)
	}

	return nil
}
