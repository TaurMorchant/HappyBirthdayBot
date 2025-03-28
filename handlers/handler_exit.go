package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	res "happy-birthday-bot/resources"
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
		msg := "Ты уверен, что хочешь уйти из программы и не желаешь быть отхеппибёзднутым?"

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Да!", okButton),
				tgbotapi.NewInlineKeyboardButtonData("ГАЛЯ, ОТМЕНА!!!", cancelButton),
			),
		)

		sentMessage := bot.SendWithPic(chatID, msg, res.Sad_cat, &inlineKeyboard, false)

		WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: h})
	} else {
		bot.SendWithPic(chatID, "Слыш, ты и так не в программе!", res.Suspicious_cat, nil, false)
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

		bot.SendWithPic(chatID, "Штош, ты удален", res.Sad_cat, nil, false)
	} else if update.CallbackQuery.Data == cancelButton {
		bot.SendWithPic(chatID, "Да ладно, ладно, не ори!", res.Do_not_scream, nil, false)
	} else {
		bot.SendWithPic(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!", res.Do_not_understand, nil, false)
	}

	return nil
}
