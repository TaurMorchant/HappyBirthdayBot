package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
	"log"
	"strings"
)

type WishlistHandler struct {
}

func (h WishlistHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()

	if user, ok := users.Get(usr.UserId(userID)); ok {
		if len(user.Wishlist) == 0 {
			msg := "Похоже ты еще не составил свой вишлист! Самое время это сделать! Напиши в ответ на это сообщение, что бы ты хотел получить в подарок?"
			bot.SendWithForceReply(chatID, msg)
			WaitForReply(usr.UserId(userID), h)
		} else {
			msg := fmt.Sprintf("Вот так выглядит твой вишлист:\n\n```\n%s\n```\n"+
				"Если хочешь его поменять, отправь мне новый вишлист в ответ на это сообщение. Если не хочешь - ответь '`Не хочу`'", user.Wishlist)
			bot.SendWithForceReply(chatID, msg)
			WaitForReply(usr.UserId(userID), h)
		}
	} else {
		bot.Send(chatID, "Кажется ты еще не зарегистрирован в программе! Зарегистрируйся при помощи команды `/join`!")
	}

	return nil
}

func (h WishlistHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	if strings.EqualFold(update.Message.Text, "не хочу") {
		bot.Send(chatID, "Океюшки")
	} else {
		users := sheets.Read()
		if user, ok := users.Get(usr.UserId(userID)); ok {
			user.Wishlist = update.Message.Text
			sheets.Write(&users)
			bot.Send(chatID, "Вишлист обновлён!")
		} else {
			log.Panicf("User with ID %d not found", usr.UserId(userID))
		}
	}
	return nil
}
