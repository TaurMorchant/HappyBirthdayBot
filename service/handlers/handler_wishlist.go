package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/mybot"
	"happy-birthday-bot/resources"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
)

type WishlistHandler struct {
}

func (h WishlistHandler) Handle(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	messageID := update.Message.MessageID

	users := sheets.Read()

	if user, ok := users.Get(usr.UserId(userID)); ok {
		if len(user.Wishlist) == 0 {
			msg := "Похоже ты еще не составил свой вишлист! Самое время это сделать! Напиши в ответ на это сообщение, что бы ты хотел получить в подарок?"
			bot.SendPicForceReply(chatID, msg, res.Wishlist, messageID)
			WaitingForReplyHandlers.Add(userID, h)
		} else {
			msg := fmt.Sprintf("У тебя сейчас такой вишлист:\n\n```\n%s\n```\n"+
				"Хочешь его поменять?", user.Wishlist)

			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Хочу", okButton),
					tgbotapi.NewInlineKeyboardButtonData("Не, все норм", cancelButton),
				),
			)
			sentMessage := bot.SendPicWithKeyboard(chatID, msg, res.Wishlist, &inlineKeyboard)
			WaitingForCallbackHandlers.Add(sentMessage.MessageID, CallbackElement{UserId: userID, Handler: h})
		}
	} else {
		bot.SendPic(chatID, "Кажется ты еще не зарегистрирован в программе! Зарегистрируйся при помощи команды [/join](/join)!", res.Suspicious)
	}

	return nil
}

func (h WishlistHandler) HandleReply(bot *mybot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	if user, ok := users.Get(usr.UserId(userID)); ok {
		user.Wishlist = update.Message.Text
		sheets.Write(&users)

		bot.SendPic(chatID, "Вжух, вишлист обновлён! 👌", res.Vjuh)
	} else {
		log.Panicf("User with ID %d not found", usr.UserId(userID))
	}
	return nil
}

func (h WishlistHandler) HandleCallback(bot *mybot.Bot, update tgbotapi.Update) error {
	log.Println("Handle callback for WishlistHandler")
	chatID := update.CallbackQuery.Message.Chat.ID
	userID := update.CallbackQuery.From.ID
	messageID := update.CallbackQuery.Message.MessageID

	if update.CallbackQuery.Data == okButton {
		msg := "Напиши в ответ на это сообщение, что бы ты хотел получить в подарок?"
		bot.SendPicForceReply(chatID, msg, res.Wishlist, messageID)
		WaitingForReplyHandlers.Add(userID, h)
	} else if update.CallbackQuery.Data == cancelButton {
		bot.SendPic(chatID, "Океюшки", res.Ok)
	} else {
		bot.SendPic(chatID, "Ты откуда вообще взял эту кнопку, тут ее не должно быть!", res.Suspicious)
	}

	return nil
}
