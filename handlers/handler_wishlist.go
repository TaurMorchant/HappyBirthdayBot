package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
)

type WishlistHandler struct {
}

func (h WishlistHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	bot.SendWithEH(tgbotapi.NewMessage(chatID, "Wishlist placeholder"))

	return nil
}

func (h WishlistHandler) HandleReply(*bot.Bot, tgbotapi.Update) error {
	return nil
}
