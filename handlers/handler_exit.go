package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"happy-birthday-bot/bot"
	"happy-birthday-bot/sheets"
	"happy-birthday-bot/usr"
)

type ExitHandler struct {
}

func (h ExitHandler) Handle(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()

	if _, ok := users.Get(usr.UserId(userID)); ok {
		msg := "Ты точно уверен, что не хочешь быть отхеппибёзднутым?\n\nЕсли уверен, ответь на  это сообщение `Да`"
		message := tgbotapi.NewMessage(chatID, msg)
		message.ParseMode = tgbotapi.ModeMarkdown
		message.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true} // Принудительный reply-режим
		bot.SendWithEH(message)

		WaitForReply(usr.UserId(userID), h)
	} else {
		bot.SendWithEH(tgbotapi.NewMessage(chatID, "Слыш, ты и так не в программе!"))
	}
	return nil
}

func (h ExitHandler) HandleReply(bot *bot.Bot, update tgbotapi.Update) error {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	users := sheets.Read()
	users.Delete(usr.UserId(userID))
	sheets.Write(&users)
	bot.SendWithEH(tgbotapi.NewMessage(chatID, "Все пучком, ты удален из программы!"))

	return nil
}
